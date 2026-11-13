// Package render contains the low-level pixel-painting primitives used by the
// galaxy package.
package render

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/amcajal/8_bit_hubble_golang/v2/internal/rngm"
	"github.com/amcajal/8_bit_hubble_golang/v2/internal/sprites"
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
// canvas. It returns a []image.Image containing one transparent layer per
// sprite size (small, medium, large, special), in that order.
//
// Each returned layer is an independent NRGBA image: compositing layers[0]
// through layers[3] onto a copy of bg exactly reconstructs the frame.
//
// The seed reset ensures that static and animated frames are visually
// consistent with each other.
func PaintForeground(canvas *image.NRGBA, width, height int, animated bool) []image.Image {
	rngm.ReloadSeed()

	layers := make([]image.Image, 4)
	layers[0] = paintSprite(canvas, sprites.Small, width, height, animated)
	layers[1] = paintSprite(canvas, sprites.Medium, width, height, animated)
	layers[2] = paintSprite(canvas, sprites.Large, width, height, animated)
	layers[3] = paintSprite(canvas, sprites.Special, width, height, animated)

	return layers
}

// paintSprite paints 0–maxLayers sprite layers of the given size onto a fresh
// transparent canvas, composites that canvas onto frame, and returns the
// transparent canvas as an independent layer image.
func paintSprite(frame *image.NRGBA, spriteSize sprites.Size, width, height int, animated bool) image.Image {
	layer := image.NewNRGBA(image.Rect(0, 0, width, height))
	// layer is zero-initialised: all pixels are fully transparent.

	if !sprites.ShouldPaintSprite(spriteSize) {
		// Return the empty transparent layer so the index contract is kept.
		return layer
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
			draw.Draw(layer, r, img, sb.Min, draw.Over)
		}
	}

	// Composite the layer onto the shared frame.
	draw.Draw(frame, frame.Bounds(), layer, image.Point{}, draw.Over)

	return layer
}

// SetNoiseSource wires the shared seeded RNG into the noise package.
// Must be called once before BuildBackground when the noise package is used.
func SetNoiseSource() {
	noise.SetRandSource(rngm.GetRand())
}
