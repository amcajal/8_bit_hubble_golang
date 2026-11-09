package render

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
	"github.com/amcajal/8_bit_hubble_golang/internal/sprites"
	"github.com/amcajal/bild_modfor8bhg/noise"
)

func BuildBackground(width, height int) image.Image {
	bg := image.NewNRGBA(image.Rect(0, 0, width, height))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(bg, bg.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	if sprites.ShouldPaintNebula() {
		nebula := sprites.GenerateNebula(width, height)
		draw.Draw(bg, bg.Bounds(), nebula, image.Point{}, draw.Over)
	}

	if sprites.ShouldPaintCosmicLatteBackground() {
		starsBackground := sprites.GetCosmicLatteBackground(width, height)
		draw.Draw(bg, bg.Bounds(), starsBackground, image.Point{}, draw.Over)
	}

	return bg
}

func PaintBackground(frame *image.NRGBA, bg image.Image) {
	draw.Draw(frame, frame.Bounds(), bg, image.Point{}, draw.Src)
}

// PaintForeground resets the RNG seed and paints all sprite layers onto the
// canvas. The seed reset ensures that static and animated frames are visually
// consistent with each other.
func PaintForeground(canvas *image.NRGBA, width, height int, animated bool) {
	rngm.ReloadSeed()
	paintSprite(canvas, sprites.Small, width, height, animated)
	paintSprite(canvas, sprites.Medium, width, height, animated)
	paintSprite(canvas, sprites.Large, width, height, animated)
	paintSprite(canvas, sprites.Special, width, height, animated)
}

// paintSprite paints 0–maxLayers layers, each with 0–maxSprites sprites of
// the given size, with a configurable probability of being skipped entirely.
func paintSprite(frame *image.NRGBA, spriteSize sprites.Size, width, height int, animated bool) {
	if !sprites.ShouldPaintSprite(spriteSize) {
		return
	}

	layers := sprites.GetNoOfLayers(spriteSize)
	for l := 1; l <= layers; l++ {
		elements := sprites.GetNoOfElements(spriteSize)
		img := sprites.GetSprite(spriteSize, animated)
		sb := img.Bounds()
		// To avoid sprites being painted off the edge of the canvas.
		widthLimit := width - sb.Max.X
		// Only width is controlled, so images in "infinite scroll" mode
		// can be scrolled horizontally without being cut off.

		for e := 1; e <= elements; e++ {
			xC := rngm.Intn(widthLimit) + 1
			yC := rngm.Intn(height)
			dp := image.Pt(xC, yC)
			r := image.Rectangle{dp, dp.Add(sb.Size())}
			draw.Draw(frame, r, img, sb.Min, draw.Over)
		}
	}
}

// SetNoiseSource wires the shared seeded RNG into the noise package.
// Must be called once before BuildBackground when the noise package is used.
func SetNoiseSource() {
	noise.SetRandSource(rngm.GetRand())
}
