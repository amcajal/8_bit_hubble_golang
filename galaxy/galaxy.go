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

	"github.com/amcajal/8_bit_hubble_golang/v2/internal/render"
	"github.com/amcajal/8_bit_hubble_golang/v2/internal/rngm"
	"github.com/amcajal/8_bit_hubble_golang/v2/internal/sprites"
)

const pngExt string = ".png"
const webpExt string = ".webp"

// ApplyConfig seeds the RNG and, when RenderConfig.RNGConfigFile is set,
// loads the TOML override file.
//
// Must be called before RenderFrames. Without it, Seed is ignored and output
// is non-deterministic. Function "Generate" calls this automatically.
func ApplyConfig(cfg RenderConfig) error {
	if cfg.RNGConfigFile != "" {
		if err := rngm.LoadConfig(cfg.RNGConfigFile); err != nil {
			return fmt.Errorf("load RNG config: %w", err)
		}
	}
	rngm.LoadSeed(cfg.Seed)
	return nil
}

// GalaxyData is returned by RenderFrames.
//
// Frames always contains the fully composited images.
// When RenderConfig.RetainLayers is true, Layers[i] holds the individual
// compositing layers for Frames[i], in order:
//
//	[0] background (black + optional nebula / cosmic-latte)
//	[1] small sprites
//	[2] medium sprites
//	[3] large sprites
//	[4] special sprites
//
// Compositing Layers[i][0..4] in order exactly reconstructs Frames[i].
// Layers is nil when RetainLayers is false.
type GalaxyData struct {
	Frames []image.Image
	Layers [][]image.Image
}

// RenderFrames generates the galaxy image data in memory.
//
//   - Animated=false -> GalaxyData.Frames has one image (use with SaveAsPng).
//   - Animated=true  -> GalaxyData.Frames has the full WebP frame sequence (use with SaveAsWebp).
//   - RetainLayers=true -> GalaxyData.Layers is populated (see GalaxyData).
//
// The background is built once and reused across all frames.
// Call ApplyConfig before this function.
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
			frameLayers := make([]image.Image, 1+len(spriteLayers))
			frameLayers[0] = bg
			copy(frameLayers[1:], spriteLayers)
			data.Layers[i] = frameLayers
		}

		sprites.AdvanceOpacityIndexes()
	}
	return data
}

// Generate is the single entry point for CLI and GUI usage.
// It validates the config, seeds the RNG, renders the frames, and saves the
// result to disk as PNG (static) or WebP (animated).
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

// SaveAsPng encodes a single image as a PNG file at the path described by cfg.
// When cfg.Paletted is true, output is converted to a 256-colour Plan 9 palette
// with Floyd-Steinberg dithering.
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

// SaveAsWebp encodes a slice of images as an animated WebP file at the path
// described by cfg. Use the slice from GalaxyData.Frames when Animated is true.
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

// appendExtension appends the correct file extension to cfg.OutputFileName
// if not already present.
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
