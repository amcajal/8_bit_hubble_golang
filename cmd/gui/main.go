package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/explorer"

	"github.com/amcajal/8_bit_hubble_golang/internal/galaxy"
	"github.com/amcajal/8_bit_hubble_golang/internal/param"
	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
)

//go:embed assets/logo.png
var assetsFS embed.FS

const appAbout = `8-Bit Hubble — Golang Version
Generate random pixel-art galaxy images!

Version: v1.0.0
Repository: https://github.com/amcajal/8_bit_hubble_golang

Original C project: https://github.com/amcajal/8_bit_hubble
Original author: Alberto Martin Cajal

8-Bit Hubble Copyright (C) 2018 Alberto Martin Cajal
Golang port: open-source, GNU GPL v3.0
This program comes with ABSOLUTELY NO WARRANTY.
For license details: https://www.gnu.org/licenses/gpl-3.0.en.html`

const appHelp = `8-Bit Hubble (Golang Version) — Quick Guide

To generate a galaxy image, fill in the form fields and press
"Generate Galaxy!":

  Output image directory
    Path to the folder where the image will be saved.
    Example: C:\Users\you\Pictures   or   ./output

  Output image name
    Base filename for the generated image (no extension needed).
    The correct extension (.png or .webp) is added automatically.
    Example: my_galaxy

  Seed (optional)
    An integer used to initialise the random number generator.
    Using the same seed always produces the same image.
    Leave blank to use the current timestamp as seed.
    Example: 42

After generation, the result image and its seed are shown in a
pop-up overlay. Copy the seed if you want to reproduce the image.

You can also run the CLI version directly:
  go run ./cmd/cli -o <output_dir> -n <name> [-s <seed>]
  go run ./cmd/cli --help`

type overlay int

const (
	overlayNone overlay = iota
	overlayHelp
	overlayAbout
	overlayResult
)

// pickerResult carries a path selected by the OS picker back to the event loop.
type pickerResult struct {
	field string // "outputDir" or "configFile"
	path  string
}

type appState struct {
	// ── Basic input widgets ───────────────────────────────────────────────
	outputDirEditor    widget.Editor
	outputDirBrowseBtn widget.Clickable
	imgNameEditor      widget.Editor
	seedEditor         widget.Editor

	// ── Advanced options (collapsed by default) ───────────────────────────
	advancedDiscloser  component.DiscloserState // collapse / expand controller
	configFileEditor   widget.Editor            // -c / --config
	configFileBrowseBtn widget.Clickable
	animatedCheck      widget.Bool              // -a / --animated
	palettedCheck      widget.Bool              // -p / --paletted
	animatedTip        component.TipArea        // hover tooltip for animated
	palettedTip        component.TipArea        // hover tooltip for paletted

	// ── Buttons ───────────────────────────────────────────────────────────
	generateBtn  widget.Clickable
	helpBtn      widget.Clickable
	aboutBtn     widget.Clickable
	closeOverlay widget.Clickable

	// ── OS file/folder picker ─────────────────────────────────────────────
	exp         *explorer.Explorer     // one instance per window
	pickerCh    chan pickerResult       // goroutines send results here
	pickerBusy  bool                   // only one dialog at a time

	// ── Scrollers ─────────────────────────────────────────────────────────
	mainList  widget.List // scrolls the entire main-window content
	helpList  widget.List
	aboutList widget.List

	// ── Status line ───────────────────────────────────────────────────────
	statusMsg string

	// ── Overlay state ─────────────────────────────────────────────────────
	activeOverlay overlay

	// ── Result overlay ────────────────────────────────────────────────────
	resultImage   image.Image
	resultImageOp paint.ImageOp
	resultSeed    int64

	// ── Logo ──────────────────────────────────────────────────────────────
	logoImage   image.Image
	logoImageOp paint.ImageOp
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("8-Bit Hubble Golang"),
			app.Size(unit.Dp(460), unit.Dp(560)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()

	state := &appState{
		statusMsg: "Ready",
		pickerCh:  make(chan pickerResult, 1),
	}
	state.mainList.Axis = layout.Vertical

	state.outputDirEditor.SingleLine = true
	state.imgNameEditor.SingleLine = true
	state.seedEditor.SingleLine = true
	state.seedEditor.Filter = "0123456789"
	state.configFileEditor.SingleLine = true

	cwd, _ := os.Getwd()
	state.outputDirEditor.SetText(cwd)
	state.imgNameEditor.SetText("8bh_galaxy")

	// Load logo
	state.logoImage = loadEmbeddedPNG("assets/logo.png")
	if state.logoImage != nil {
		state.logoImageOp = paint.NewImageOp(state.logoImage)
	}

	// Create the OS file/folder explorer (one per window, must be created before the loop).
	state.exp = explorer.NewExplorer(w)

	var ops op.Ops

	for {
		e := w.Event()

		// The explorer needs to see every event to work correctly on all platforms.
		state.exp.ListenEvents(e)

		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Drain any pending picker results (sent from goroutines).
			select {
			case res := <-state.pickerCh:
				state.pickerBusy = false
				switch res.field {
				case "outputDir":
					state.outputDirEditor.SetText(res.path)
				case "configFile":
					state.configFileEditor.SetText(res.path)
				}
				w.Invalidate()
			default:
			}

			// Handle button clicks
			if state.generateBtn.Clicked(gtx) {
				state.handleGenerate(w)
			}
			if state.helpBtn.Clicked(gtx) {
				state.activeOverlay = overlayHelp
			}
			if state.aboutBtn.Clicked(gtx) {
				state.activeOverlay = overlayAbout
			}
			if state.closeOverlay.Clicked(gtx) {
				state.activeOverlay = overlayNone
			}

			// Browse buttons — launch OS dialogs in goroutines.
			if !state.pickerBusy && state.outputDirBrowseBtn.Clicked(gtx) {
				state.pickerBusy = true
				go func() {
					path, err := state.exp.ChooseFolder()
					if err != nil {
						state.pickerBusy = false
						return
					}
					state.pickerCh <- pickerResult{field: "outputDir", path: path}
					w.Invalidate()
				}()
			}
			if !state.pickerBusy && state.configFileBrowseBtn.Clicked(gtx) {
				state.pickerBusy = true
				go func() {
					f, err := state.exp.ChooseFile(".toml")
					if err != nil {
						state.pickerBusy = false
						return
					}
					defer f.Close()
					// ChooseFile returns an io.ReadCloser; on desktop it is an *os.File.
					type namer interface{ Name() string }
					path := ""
					if named, ok := f.(namer); ok {
						path = filepath.Clean(named.Name())
					}
					if path != "" {
						state.pickerCh <- pickerResult{field: "configFile", path: path}
						w.Invalidate()
					} else {
						state.pickerBusy = false
					}
				}()
			}

			// Draw
			drawUI(gtx, th, state)

			e.Frame(gtx.Ops)
		}
	}
}

func (s *appState) handleGenerate(w *app.Window) {
	outputDir := s.outputDirEditor.Text()
	imgName := s.imgNameEditor.Text()
	seedStr := s.seedEditor.Text()

	if outputDir == "" {
		s.statusMsg = "Error: output directory is empty."
		return
	}
	if imgName == "" {
		s.statusMsg = "Error: image name is empty."
		return
	}

	// Parse the seed — default to current timestamp when blank.
	var seed int64
	if seedStr != "" {
		if _, err := fmt.Sscanf(seedStr, "%d", &seed); err != nil {
			s.statusMsg = "Error: seed must be an integer."
			return
		}
	} else {
		seed = time.Now().Unix()
	}

	if err := param.CheckParamsDirect(
		outputDir,
		imgName,
		s.configFileEditor.Text(),
		seed,
		s.animatedCheck.Value,
		s.palettedCheck.Value,
	); err != nil {
		s.statusMsg = fmt.Sprintf("Parameter error: %v", err)
		return
	}

	rngm.LoadSeed(param.InputParams.Seed)

	if err := galaxy.GenerateGalaxy(); err != nil {
		s.statusMsg = fmt.Sprintf("Generation error: %v", err)
		return
	}

	s.statusMsg = fmt.Sprintf("Galaxy generated! Seed: %d", param.InputParams.Seed)
	s.resultSeed = param.InputParams.Seed

	// Load the generated image for the result overlay
	fullPath := outputDir + "/" + param.InputParams.OutputFileName
	f, err := os.Open(fullPath)
	if err == nil {
		defer f.Close()
		img, _, err := image.Decode(f)
		if err == nil {
			s.resultImage = img
			s.resultImageOp = paint.NewImageOp(img)
			s.activeOverlay = overlayResult
		}
	}

	w.Invalidate()
}

var bgColor = color.NRGBA{R: 0xBE, G: 0xBE, B: 0xBE, A: 0xFF}
var accentColor = color.NRGBA{R: 0x4B, G: 0x69, B: 0x83, A: 0xFF}

func drawUI(gtx layout.Context, th *material.Theme, s *appState) layout.Dimensions {
	paint.Fill(gtx.Ops, bgColor)

	// All main-window content lives inside a vertical List so it scrolls when
	// the window is too short (e.g. when Advanced Options are expanded).
	dims := material.List(th, &s.mainList).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				// ── Logo ──────────────────────────────────────────────────────
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return drawLogo(gtx, s)
					})
				}),

				layout.Rigid(vspacer(8)),

				// ── Output directory ──────────────────────────────────────────
				layout.Rigid(boldLabel(th, "Output image directory:")),
				layout.Rigid(vspacer(4)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawPathField(gtx, th, &s.outputDirEditor, "e.g. ./output",
						&s.outputDirBrowseBtn, s.pickerBusy)
				}),

				layout.Rigid(vspacer(8)),

				// ── Image name ────────────────────────────────────────────────
				layout.Rigid(boldLabel(th, "Output image name:")),
				layout.Rigid(vspacer(4)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawEditor(gtx, th, &s.imgNameEditor, "e.g. my_galaxy")
				}),

				layout.Rigid(vspacer(8)),

				// ── Seed ──────────────────────────────────────────────────────
				layout.Rigid(boldLabel(th, "Seed value (optional):")),
				layout.Rigid(vspacer(4)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawEditor(gtx, th, &s.seedEditor, "e.g. 42")
				}),

				layout.Rigid(vspacer(10)),

				// ── Advanced options (collapsible) ────────────────────────────
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawAdvancedOptions(gtx, th, s)
				}),

				layout.Rigid(vspacer(10)),

				// ── Generate button ───────────────────────────────────────────
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &s.generateBtn, "Generate Galaxy!")
					btn.Background = accentColor
					return btn.Layout(gtx)
				}),

				layout.Rigid(vspacer(8)),

				// ── Status console ────────────────────────────────────────────
				layout.Rigid(boldLabel(th, "Status console:")),
				layout.Rigid(vspacer(4)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawStatusBox(gtx, th, s.statusMsg)
				}),

				layout.Rigid(vspacer(12)),

				// ── Help / About buttons ──────────────────────────────────────
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &s.helpBtn, "Help")
								btn.Background = accentColor
								return btn.Layout(gtx)
							})
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &s.aboutBtn, "About")
								btn.Background = accentColor
								return btn.Layout(gtx)
							})
						}),
					)
				}),
			)
		})
	})

	switch s.activeOverlay {
	case overlayHelp:
		drawTextOverlay(gtx, th, s, "Help", appHelp, &s.helpList)
	case overlayAbout:
		drawTextOverlay(gtx, th, s, "About 8-Bit Hubble", appAbout, &s.aboutList)
	case overlayResult:
		drawResultOverlay(gtx, th, s)
	}

	return dims
}

func drawAdvancedOptions(gtx layout.Context, th *material.Theme, s *appState) layout.Dimensions {
	sd := component.SimpleDiscloser(th, &s.advancedDiscloser)

	dividerColor := color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}

	return widget.Border{
		Color:        dividerColor,
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(6),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return sd.Layout(gtx,
				// summary header
				func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, "Advanced Options")
					lbl.Font.Weight = 700
					lbl.Color = accentColor
					return lbl.Layout(gtx)
				},
				// detail body
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(vspacer(8)),

						// Config file
						layout.Rigid(boldLabel(th, "Config file (optional):")),
						layout.Rigid(vspacer(4)),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return drawPathField(gtx, th, &s.configFileEditor, "path/to/config.toml",
								&s.configFileBrowseBtn, s.pickerBusy)
						}),

						layout.Rigid(vspacer(10)),

						// Animated checkbox with tooltip
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							tip := component.DesktopTooltip(th,
								"Generate a WebP animated image instead of a static PNG.\nSprites will blink by cycling through opacity levels.")
							return s.animatedTip.Layout(gtx, tip, func(gtx layout.Context) layout.Dimensions {
								cb := material.CheckBox(th, &s.animatedCheck, "Animated (WebP output)")
								return cb.Layout(gtx)
							})
						}),

						layout.Rigid(vspacer(6)),

						// Paletted checkbox with tooltip
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							tip := component.DesktopTooltip(th,
								"Restrict the output image to 256 colours\nfrom the Plan 9 palette (Floyd-Steinberg dithering).")
							return s.palettedTip.Layout(gtx, tip, func(gtx layout.Context) layout.Dimensions {
								cb := material.CheckBox(th, &s.palettedCheck, "Paletted (256-colour PNG)")
								return cb.Layout(gtx)
							})
						}),

						layout.Rigid(vspacer(4)),
					)
				},
			)
		})
	})
}

func drawLogo(gtx layout.Context, s *appState) layout.Dimensions {
	if s.logoImage == nil {
		return layout.Dimensions{}
	}
	b := s.logoImage.Bounds()
	// Scale so the logo fits within 200 dp height, preserving aspect ratio.
	maxH := float32(gtx.Dp(200))
	scale := maxH / float32(b.Dy())
	imgWidget := widget.Image{
		Src:   s.logoImageOp,
		Fit:   widget.Contain,
		Scale: scale / (float32(gtx.Metric.PxPerDp)),
	}
	return imgWidget.Layout(gtx)
}

// drawPathField renders a single-line editor + a small "Browse…" button on the right.
// The button is disabled (greyed out) while another picker dialog is already open.
func drawPathField(gtx layout.Context, th *material.Theme, ed *widget.Editor, hint string,
	browseBtn *widget.Clickable, busy bool) layout.Dimensions {

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return drawEditor(gtx, th, ed, hint)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Dim the button while a dialog is already open.
			btn := material.Button(th, browseBtn, "Browse…")
			if busy {
				btn.Background = color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF}
			} else {
				btn.Background = accentColor
			}
			return btn.Layout(gtx)
		}),
	)
}

func drawEditor(gtx layout.Context, th *material.Theme, ed *widget.Editor, hint string) layout.Dimensions {
	return widget.Border{
		Color:        color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(6)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(th, ed, hint)
			e.Color = color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF}
			return e.Layout(gtx)
		})
	})
}

func drawStatusBox(gtx layout.Context, th *material.Theme, msg string) layout.Dimensions {
	return widget.Border{
		Color:        color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Give the status box a fixed minimum height
		gtx.Constraints.Min.Y = gtx.Dp(48)
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(th, msg)
			lbl.Color = color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF}
			lbl.WrapPolicy = text.WrapWords
			return lbl.Layout(gtx)
		})
	})
}

// drawTextOverlay renders a semi-transparent modal panel with scrollable text.
func drawTextOverlay(gtx layout.Context, th *material.Theme, s *appState, title, body string, list *widget.List) {
	// Dim the background
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{A: 0xCC})

	// Panel occupies most of the window height so text has room to breathe.
	panelW := gtx.Dp(380)
	panelH := gtx.Constraints.Max.Y - gtx.Dp(60) // 30 dp margin top + bottom
	if panelH < gtx.Dp(200) {
		panelH = gtx.Dp(200)
	}
	ox := (gtx.Constraints.Max.X - panelW) / 2
	oy := (gtx.Constraints.Max.Y - panelH) / 2

	defer op.Offset(image.Pt(ox, oy)).Push(gtx.Ops).Pop()
	gtx.Constraints = layout.Exact(image.Pt(panelW, panelH))

	drawPanel(gtx, th, s, title, body, list)
}

func drawPanel(gtx layout.Context, th *material.Theme, s *appState, title, body string, list *widget.List) {
	// Panel background
	rr := clip.RRect{
		Rect: image.Rectangle{Max: gtx.Constraints.Max},
		NE:   8, NW: 8, SE: 8, SW: 8,
	}
	defer rr.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF5, A: 0xFF})

	layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Title row + close button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.H6(th, title)
						lbl.Color = accentColor
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &s.closeOverlay, "✕")
						btn.Background = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
						return btn.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(vspacer(8)),
			// Scrollable body — the list has a single item (the full text block)
			// so scroll state is tracked by widget.List.
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				list.Axis = layout.Vertical
				return material.List(th, list).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
					lbl := material.Body2(th, body)
					lbl.Color = color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
					lbl.WrapPolicy = text.WrapWords
					return lbl.Layout(gtx)
				})
			}),
		)
	})
}

// drawResultOverlay shows the generated galaxy image with its seed.
func drawResultOverlay(gtx layout.Context, th *material.Theme, s *appState) {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{A: 0xCC})

	panelW := gtx.Dp(440)
	panelH := gtx.Dp(480)
	ox := (gtx.Constraints.Max.X - panelW) / 2
	oy := (gtx.Constraints.Max.Y - panelH) / 2

	defer op.Offset(image.Pt(ox, oy)).Push(gtx.Ops).Pop()
	gtx.Constraints = layout.Exact(image.Pt(panelW, panelH))

	rr := clip.RRect{
		Rect: image.Rectangle{Max: gtx.Constraints.Max},
		NE:   8, NW: 8, SE: 8, SW: 8,
	}
	defer rr.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF5, A: 0xFF})

	layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Title + close
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.H6(th, "Result")
						lbl.Color = accentColor
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &s.closeOverlay, "✕")
						btn.Background = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
						return btn.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(vspacer(8)),
			// Image
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if s.resultImage == nil {
					return layout.Dimensions{}
				}
				b := s.resultImage.Bounds()
				// Fit within the available space
				availW := float32(gtx.Constraints.Max.X)
				availH := float32(gtx.Constraints.Max.Y)
				scaleX := availW / float32(b.Dx())
				scaleY := availH / float32(b.Dy())
				scale := scaleX
				if scaleY < scale {
					scale = scaleY
				}
				imgWidget := widget.Image{
					Src:      s.resultImageOp,
					Fit:      widget.Contain,
					Position: layout.Center,
					Scale:    scale / float32(gtx.Metric.PxPerDp),
				}
				return imgWidget.Layout(gtx)
			}),
			layout.Rigid(vspacer(8)),
			// Seed label (selectable text)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, fmt.Sprintf("SEED OF THE GALAXY: %d", s.resultSeed))
				lbl.Color = color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF}
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			}),
		)
	})
}

func boldLabel(th *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body1(th, txt)
		lbl.Font.Weight = 700 // semibold
		lbl.Color = color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF}
		return lbl.Layout(gtx)
	}
}

func vspacer(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: dp}.Layout(gtx)
	}
}

func loadEmbeddedPNG(path string) image.Image {
	f, err := assetsFS.Open(path)
	if err != nil {
		log.Printf("warning: could not open embedded asset %q: %v", path, err)
		return nil
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Printf("warning: could not decode embedded asset %q: %v", path, err)
		return nil
	}
	return img
}
