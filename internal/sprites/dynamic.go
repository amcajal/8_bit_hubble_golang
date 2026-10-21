package sprites

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/noise"
)

func blackToTransparent(img image.Image) *image.RGBA {
	return adjust.Apply(img, func(c color.RGBA) color.RGBA {
		if c.R == 0 && c.G == 0 && c.B == 0 {
			c.A = 0 // Transparente
		}
		return c
	})
}

func GetCosmicLatteBackground() image.Image {
	cosmic_latte := color.RGBA{R: 255, G: 248, B: 231, A: 255}
	aux_img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(aux_img, aux_img.Bounds(), &image.Uniform{cosmic_latte}, image.Point{}, draw.Src)

	stars_background := noise.Generate(500, 500, &noise.Options{Monochrome: true, NoiseFn: noise.Binary})
	for i := 0; i < rngm.Intn(3)+4; i++ {
		aux_background := noise.Generate(500, 500, &noise.Options{Monochrome: true, NoiseFn: noise.Binary})
		stars_background = blend.Multiply(stars_background, aux_background)
	}
	stars_background = blend.Multiply(stars_background, aux_img)
	stars_background = blackToTransparent(stars_background)
	return stars_background
}

func GenerateColoredPerlinNoise(baseColor *color.RGBA) image.Image {
	color_to_use := color.RGBA{R: 255, G: 248, B: 231, A: 255} // Cosmic Latte
	if baseColor != nil {
		color_to_use = *baseColor
	}
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(img, img.Bounds(), &image.Uniform{
		C: color_to_use,
	}, image.Point{}, draw.Src)
	perlin := noise.GeneratePerlin(500, 500, 0.15)
	return blend.Multiply(perlin, img)
}

func GenerateNebula() image.Image {
	nebula := GenerateColoredPerlinNoise(nil)

	for i := 0; i < rngm.Intn(3)+1; i++ {
		aux_color := color.RGBA{
			R: uint8(rngm.Intn(256)),
			G: uint8(rngm.Intn(256)),
			B: uint8(rngm.Intn(256)),
			A: 255,
		}
		aux_nebula := GenerateColoredPerlinNoise(&aux_color)
		nebula = blend.Overlay(nebula, aux_nebula)
	}
	return nebula
}
