package galaxy

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/HugoSmits86/nativewebp"

	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
	"github.com/amcajal/8_bit_hubble_golang/internal/sprites"
	"github.com/amcajal/bild_modfor8bhg/noise"
)

const pngExt string = ".png"
const webpExt string = ".webp"

// GenerateImage returns a single static galaxy image in memory.
// For animated output, use GenerateToFile.
func GenerateImage(cfg Config) (image.Image, error) {
	if cfg.Animated {
		return nil, fmt.Errorf("animated output requires GenerateToFile")
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if err := applyConfig(cfg); err != nil {
		return nil, err
	}
	return generateSingleImage(cfg), nil
}

// GenerateToFile writes a galaxy image (PNG or animated WebP) to disk.
func GenerateToFile(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := applyConfig(cfg); err != nil {
		return err
	}
	appendExtension(&cfg)
	if cfg.Animated {
		return saveAsWebp(cfg, generateFrames(cfg))
	}
	return saveAsPng(cfg, generateSingleImage(cfg))
}

// applyConfig loads the optional TOML config file and seeds the RNG.
func applyConfig(cfg Config) error {
	if cfg.ConfigFile != "" {
		if err := rngm.LoadConfig(cfg.ConfigFile); err != nil {
			return fmt.Errorf("load config: %w", err)
		}
	}
	rngm.LoadSeed(cfg.Seed)
	return nil
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

func generateSingleImage(cfg Config) image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
	noise.SetRandSource(rngm.GetRand())
	paintBackground(canvas, buildBackground(cfg), cfg)
	paintForeground(canvas, cfg)
	return canvas
}

func generateFrames(cfg Config) []image.Image {
	noOfFrames := len(sprites.OpacityValues)
	images := make([]image.Image, noOfFrames)

	noise.SetRandSource(rngm.GetRand())
	bg := buildBackground(cfg)

	sprites.InitOpacityIndexes()
	for i := 0; i < noOfFrames; i++ {
		frame := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
		paintBackground(frame, bg, cfg)
		paintForeground(frame, cfg)
		images[i] = frame
		sprites.AdvanceOpacityIndexes()
	}
	return images
}

func buildBackground(cfg Config) image.Image {
	bg := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(bg, bg.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	if sprites.ShouldPaintNebula() {
		nebula := sprites.GenerateNebula(cfg.Width, cfg.Height)
		draw.Draw(bg, bg.Bounds(), nebula, image.Point{}, draw.Over)
	}

	if sprites.ShouldPaintCosmicLatteBackground() {
		stars_background := sprites.GetCosmicLatteBackground(cfg.Width, cfg.Height)
		draw.Draw(bg, bg.Bounds(), stars_background, image.Point{}, draw.Over)
	}

	return bg
}

func paintBackground(frame *image.NRGBA, bg image.Image, cfg Config) {
	draw.Draw(frame, frame.Bounds(), bg, image.Point{}, draw.Src)
}

func paintForeground(canvas *image.NRGBA, cfg Config) {
	rngm.ReloadSeed() // Reset seed so animated and static frames match
	paintSprite(canvas, sprites.Small, cfg)
	paintSprite(canvas, sprites.Medium, cfg)
	paintSprite(canvas, sprites.Large, cfg)
	paintSprite(canvas, sprites.Special, cfg)
}

// paintSprite paints 0–maxLayers layers, each with 0–maxSprites sprites of
// the given size, with a configurable probability of being skipped entirely.
func paintSprite(frame *image.NRGBA, spriteSize sprites.Size, cfg Config) {
	if !sprites.ShouldPaintSprite(spriteSize) {
		return
	}

	layers := sprites.GetNoOfLayers(spriteSize)
	for l := 1; l <= layers; l++ {
		elements := sprites.GetNoOfElements(spriteSize)
		img := sprites.GetSprite(spriteSize, cfg.Animated)
		sb := img.Bounds()
		// To avoid sprites being painted off the edge of the canvas
		width_limit := cfg.Width - sb.Max.X
		// Only width is controled, so images in "infinite scroll" mode
		// can be scrolled horizontally without being cut off.

		for e := 1; e <= elements; e++ {
			x_c := rngm.Intn(width_limit) + 1
			y_c := rngm.Intn(cfg.Height)
			dp := image.Pt(x_c, y_c)
			r := image.Rectangle{dp, dp.Add(sb.Size())}
			draw.Draw(frame, r, img, sb.Min, draw.Over)
		}
	}
}

func saveAsPng(cfg Config, frame image.Image) error {
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

func saveAsWebp(cfg Config, images []image.Image) error {
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
