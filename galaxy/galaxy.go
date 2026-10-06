// Package galaxy creates the "8-bit galaxy" image (sprites, colors, etc)
package galaxy

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/amcajal/8_bit_hubble_golang/palette"
	"github.com/amcajal/8_bit_hubble_golang/param"
	"github.com/amcajal/8_bit_hubble_golang/rngm"
	"github.com/amcajal/8_bit_hubble_golang/sprites"
)

// Output image (the canvas)
var canvas *image.NRGBA

// Fixed dimensions of the output image
const dim_x int = 500 // width
const dim_y int = 500 // height

func GenerateGalaxy() error {

	// Set canvas dimensions
	canvas = image.NewNRGBA(image.Rect(0, 0, dim_x, dim_y))

	// Paint background
	paintBackground()

	// Paint small stars
	paintSprite(sprites.Small, rngm.Cfg.SmallSpriteMaxSprites)

	// Paint medium stars
	paintSprite(sprites.Medium, rngm.Cfg.MediumSpriteMaxSprites)

	// Paint big stars
	paintSprite(sprites.Large, rngm.Cfg.LargeSpriteMaxSprites)

	// Paint special sprites
	paintSprite(sprites.Special, rngm.Cfg.SpecialSpriteMaxSprites)

	// Save image
	writer, err := os.Create(param.OutputDir + "/" + param.PngName)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err = png.Encode(writer, canvas); err != nil {
		return err
	}

	return nil
}

func paintBackground() {
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)
}

// There is a chance% probability of painting sprites. If success, paint
// 0 to maxLayers, each one with 0 to maxSprites of spriteSize type.
func paintSprite(spriteSize sprites.Size, maxSprites int) {

	if !sprites.ShouldPaintSprite(spriteSize) {
		return
	}

	// Number of layers for this sprite size
	layers := sprites.GetNoOfLayers(spriteSize)

	for l := 1; l <= layers; l++ {

		// Number of sprites to be painted
		elements := rngm.Intn(maxSprites + 1)

		// Decide sprite to be painted
		sprite := sprites.GetSprite(spriteSize)

		// sb stands for "sprite bounds"
		sb := sprite.Bounds()

		// Decide color
		changeColor(&sprite)

		for e := 1; e <= elements; e++ {

			// Random position (coordinates) in the image
			x_c := rngm.Intn(dim_x)
			y_c := rngm.Intn(dim_y)

			// Draw the sprite
			dp := image.Pt(x_c, y_c)
			r := image.Rectangle{dp, dp.Add(sb.Size())}
			draw.Draw(canvas, r, sprite, sb.Min, draw.Over)
		}
	}
}

func changeColor(sprite *image.Image) {

	// sb stands for "sprite bounds"
	sb := (*sprite).Bounds()

	// Change the hue of the sprite
	palette.SetHueRotation(rngm.Intn(361)) // 0 degrees to 360 degrees

	// Colorize pixels of the sprite
	rows, columns := sb.Max.X, sb.Max.Y
	for col := 0; col < columns; col++ {
		for row := 0; row < rows; row++ {
			sprite_color := (*sprite).At(row, col)
			r, g, b, a := sprite_color.RGBA()
			if a != 0 {

				nr, ng, nb := palette.ChangeHue(r, g, b)

				var new_rgb color.RGBA = color.RGBA{
					uint8(nr & 0xFF),
					uint8(ng & 0xFF),
					uint8(nb & 0xFF),
					uint8(a & 0xFF)}

				// All standard sprites should be of thesame type, but for some
				// reason, program detects some of them as NRGBA, and other as RGBA
				switch (*sprite).(type) {
				case *image.RGBA:
					(*sprite).(*image.RGBA).Set(row, col, new_rgb)
				case *image.NRGBA:
					(*sprite).(*image.NRGBA).Set(row, col, new_rgb)
				default:
					panic("Error: Invalid Sprite type (neither RGBA nor NRGBA)")
				}

			}
		}
	}
}
