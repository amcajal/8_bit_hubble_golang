package galaxy

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/HugoSmits86/nativewebp"

	"github.com/amcajal/8_bit_hubble_golang/internal/render"
	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
	"github.com/amcajal/8_bit_hubble_golang/internal/sprites"
)

const pngExt string = ".png"
const webpExt string = ".webp"

// ApplyConfig seeds the RNG with cfg.Seed and, when cfg.RNGConfigFile is
// non-empty, loads the TOML file that overrides sprite probability parameters.
//
// Library consumers MUST call ApplyConfig before RenderFrames. Without it the
// RNG is in an undefined state and cfg.Seed is silently ignored, producing
// non-deterministic output.
//
// Generate calls ApplyConfig automatically; CLI/GUI callers do not need to
// call it directly.
func ApplyConfig(cfg RenderConfig) error {
	if cfg.RNGConfigFile != "" {
		if err := rngm.LoadConfig(cfg.RNGConfigFile); err != nil {
			return fmt.Errorf("load RNG config: %w", err)
		}
	}
	rngm.LoadSeed(cfg.Seed)
	return nil
}

// GalaxyData is the value returned by RenderFrames. It always contains the
// fully composited frames in Frames. When RenderConfig.RetainLayers is true,
// Layers[i] holds the individual compositing layers for Frames[i], in order:
//
//	[0] background (black fill + optional nebula/cosmic-latte)
//	[1] small sprites
//	[2] medium sprites
//	[3] large sprites
//	[4] special sprites
//
// Compositing Layers[i][0] through Layers[i][4] in order exactly reconstructs
// Frames[i]. When RetainLayers is false, Layers is nil.
type GalaxyData struct {
	Frames []image.Image
	Layers [][]image.Image
}

func RenderFrames(cfg RenderConfig) GalaxyData {
	render.SetNoiseSource()
	bg := render.BuildBackground(cfg.Width, cfg.Height)

	noOfFrames := 1
	if cfg.Animated {
		noOfFrames = len(sprites.OpacityValues)
		sprites.InitOpacityIndexes()
	}

	data := GalaxyData{
		Frames: make([]image.Image, noOfFrames),
	}
	if cfg.RetainLayers {
		data.Layers = make([][]image.Image, noOfFrames)
	}

	for i := 0; i < noOfFrames; i++ {
		frame := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
		render.PaintBackground(frame, bg)
		spriteLayers := render.PaintForeground(frame, cfg.Width, cfg.Height, cfg.Animated)
		data.Frames[i] = frame

		if cfg.RetainLayers {
			// Layer 0 is the background; layers 1-4 are the sprite layers
			// returned by PaintForeground (small, medium, large, special).
			frameLayers := make([]image.Image, 1+len(spriteLayers))
			frameLayers[0] = bg
			copy(frameLayers[1:], spriteLayers)
			data.Layers[i] = frameLayers
		}

		sprites.AdvanceOpacityIndexes()
	}
	return data
}

// Generate is the all-in-one entry point for the CLI and GUI paths.
// It validates the config, seeds the RNG, renders the frames, and writes the
// result to disk as a PNG (static) or WebP (animated).
//
// Library consumers should use ApplyConfig + RenderFrames + SaveAsPng /
// SaveAsWebp directly for finer control.
func Generate(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := ApplyConfig(cfg.RenderConfig); err != nil {
		return err
	}
	appendExtension(&cfg)

	data := RenderFrames(cfg.RenderConfig)

	if cfg.Animated {
		return SaveAsWebp(cfg.OutputConfig, data.Frames)
	}
	return SaveAsPng(cfg.OutputConfig, data.Frames[0])
}

func SaveAsPng(cfg OutputConfig, frame image.Image) error {
	writer, err := os.Create(cfg.OutputDir + "/" + cfg.OutputFileName)
	if err != nil {
		return err
	}
	defer writer.Close()

	var outputImage image.Image = frame
	if cfg.Paletted {
		palettedImg := image.NewPaletted(frame.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(palettedImg, palettedImg.Bounds(), frame, frame.Bounds().Min)
		outputImage = palettedImg
	}
	return png.Encode(writer, outputImage)
}

func SaveAsWebp(cfg OutputConfig, images []image.Image) error {
	writer, err := os.Create(cfg.OutputDir + "/" + cfg.OutputFileName)
	if err != nil {
		return err
	}
	defer writer.Close()

	durations := make([]uint, len(images))
	disposals := make([]uint, len(images))
	for i := range images {
		durations[i] = 175
		disposals[i] = 0
	}

	ani := nativewebp.Animation{
		Images:          images,
		Durations:       durations,
		Disposals:       disposals,
		LoopCount:       0,
		BackgroundColor: 0xffffffff,
	}
	return nativewebp.EncodeAll(writer, &ani, nil)
}

func appendExtension(cfg *Config) {
	if cfg.Animated {
		if !strings.HasSuffix(cfg.OutputFileName, webpExt) {
			cfg.OutputFileName += webpExt
		}
	} else {
		if !strings.HasSuffix(cfg.OutputFileName, pngExt) {
			cfg.OutputFileName += pngExt
		}
	}
}
