package sprites

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/amcajal/8_bit_hubble_golang/v2/internal/rngm"
	"github.com/amcajal/bild_modfor8bhg/adjust"
	"github.com/amcajal/bild_modfor8bhg/blend"
	"github.com/amcajal/bild_modfor8bhg/noise"
)

func blackToTransparent(img image.Image) *image.RGBA {
	return adjust.Apply(img, func(c color.RGBA) color.RGBA {
		if c.R == 0 && c.G == 0 && c.B == 0 {
			c.A = 0 // Transparent
		}
		return c
	})
}

func GetCosmicLatteBackground(w, h int) image.Image {
	cosmic_latte := color.RGBA{R: 255, G: 248, B: 231, A: 255}
	aux_img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(aux_img, aux_img.Bounds(), &image.Uniform{cosmic_latte}, image.Point{}, draw.Src)

	stars_background := noise.Generate(w, h, &noise.Options{Monochrome: true, NoiseFn: noise.Binary})
	for i := 0; i < rngm.Intn(3)+4; i++ {
		aux_background := noise.Generate(w, h, &noise.Options{Monochrome: true, NoiseFn: noise.Binary})
		stars_background = blend.Multiply(stars_background, aux_background)
	}
	stars_background = blend.Multiply(stars_background, aux_img)
	stars_background = blackToTransparent(stars_background)
	return stars_background
}

func GenerateColoredPerlinNoise(w, h int, baseColor *color.RGBA) image.Image {
	color_to_use := color.RGBA{R: 255, G: 248, B: 231, A: 255} // Cosmic Latte
	if baseColor != nil {
		color_to_use = *baseColor
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{
		C: color_to_use,
	}, image.Point{}, draw.Src)
	perlin := noise.GeneratePerlin(w, h, 0.15)
	return blend.Multiply(perlin, img)
}

func GenerateNebula(w, h int, seamless bool) image.Image {
	nebula := GenerateColoredPerlinNoise(w, h, nil)

	for i := 0; i < rngm.Intn(3)+1; i++ {
		aux_color := color.RGBA{
			R: uint8(rngm.Intn(256)),
			G: uint8(rngm.Intn(256)),
			B: uint8(rngm.Intn(256)),
			A: 255,
		}
		aux_nebula := GenerateColoredPerlinNoise(w, h, &aux_color)
		nebula = blend.Overlay(nebula, aux_nebula)
	}

	if seamless {
		nebula = applyNebulaSeam(nebula, w, h)
	}
	return nebula
}

// applyNebulaSeam returns a new image equal to src with a seam overlay on its
// left edge. The overlay is built from the rightmost seamWidth columns of src:
//
//  1. Those columns are mirrored horizontally (column order reversed).
//  2. A left-opaque → right-transparent alpha gradient is applied so the
//     overlay blends smoothly into the existing left edge rather than
//     replacing it abruptly.
//  3. The result is composited over src with draw.Over.
//
// seamWidth = w/10, clamped to at least 1.
func applyNebulaSeam(src image.Image, w, h int) image.Image {
	seamWidth := w / 10
	if seamWidth < 1 {
		seamWidth = 1
	}

	// Build the mirrored+faded overlay strip (seamWidth × h).
	overlay := image.NewNRGBA(image.Rect(0, 0, seamWidth, h))
	bounds := src.Bounds()
	for x := 0; x < seamWidth; x++ {
		// Mirror: x=0 maps to the rightmost column of src, x=seamWidth-1 to
		// the column at w-seamWidth (the left edge of the cut region).
		srcX := bounds.Min.X + (w - 1 - x)

		// Alpha gradient: fully opaque at x=0 (left of pasted strip) fading
		// to fully transparent at x=seamWidth-1 (right of pasted strip).
		var a uint8
		if seamWidth > 1 {
			a = uint8((seamWidth - 1 - x) * 255 / (seamWidth - 1))
		} else {
			a = 255
		}

		for y := 0; y < h; y++ {
			r32, g32, b32, _ := src.At(srcX, bounds.Min.Y+y).RGBA()
			overlay.SetNRGBA(x, y, color.NRGBA{
				R: uint8(r32 >> 8),
				G: uint8(g32 >> 8),
				B: uint8(b32 >> 8),
				A: a,
			})
		}
	}

	// Composite: copy src into dst, then draw overlay over its left columns.
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	seamRect := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+seamWidth, bounds.Max.Y)
	draw.Draw(dst, seamRect, overlay, image.Point{0, 0}, draw.Over)
	return dst
}
