// Package galaxy creates the "8-bit galaxy" image (sprites, colors, etc)
package galaxy

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/png"
	"os"

	"github.com/HugoSmits86/nativewebp"

	hubble_palette "github.com/amcajal/8_bit_hubble_golang/palette"
	"github.com/amcajal/8_bit_hubble_golang/param"
	"github.com/amcajal/8_bit_hubble_golang/rngm"
	"github.com/amcajal/8_bit_hubble_golang/sprites"
)

// Fixed dimensions of the output image
const dim_x int = 500 // width
const dim_y int = 500 // height

func GenerateGalaxy() error {
	var noOfFrames int = 1
	if param.InputParams.Animated {
		noOfFrames = 10
	}
	images := make([]image.Image, noOfFrames)

	for i := 0; i < noOfFrames; i++ {
		frame := image.NewNRGBA(image.Rect(0, 0, dim_x, dim_y))
		images[i] = frame
		paintBackground(frame)
		paintSprite(frame, sprites.Small)
		paintSprite(frame, sprites.Medium)
		paintSprite(frame, sprites.Large)
		paintSprite(frame, sprites.Special)
	}

	if param.InputParams.Animated {
		return saveAsWebp(images)
	} else {
		return saveAsPng(images[0])
	}
}

func paintBackground(frame *image.NRGBA) {
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(frame, frame.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	if sprites.ShouldPaintNebula() {
		nebula := sprites.GenerateNebula()
		draw.Draw(frame, frame.Bounds(), nebula, image.Point{}, draw.Over)
	}

	if sprites.ShouldPaintCosmicLatteBackground() {
		stars_background := sprites.GetCosmicLatteBackground()
		draw.Draw(frame, frame.Bounds(), stars_background, image.Point{}, draw.Over)
	}
}

// There is a chance% probability of painting sprites. If success, paint
// 0 to maxLayers, each one with 0 to maxSprites of spriteSize type.
func paintSprite(frame *image.NRGBA, spriteSize sprites.Size) {

	if !sprites.ShouldPaintSprite(spriteSize) {
		return
	}

	// Number of layers for this sprite size
	layers := sprites.GetNoOfLayers(spriteSize)

	for l := 1; l <= layers; l++ {

		// Number of sprites to be painted
		elements := sprites.GetNoOfElements(spriteSize)

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
			draw.Draw(frame, r, sprite, sb.Min, draw.Over)
		}
	}
}

func changeColor(sprite *image.Image) {

	// sb stands for "sprite bounds"
	sb := (*sprite).Bounds()

	// Change the hue of the sprite
	hubble_palette.SetHueRotation(rngm.Intn(361)) // 0 degrees to 360 degrees

	// Colorize pixels of the sprite
	rows, columns := sb.Max.X, sb.Max.Y
	for col := 0; col < columns; col++ {
		for row := 0; row < rows; row++ {
			sprite_color := (*sprite).At(row, col)
			r, g, b, a := sprite_color.RGBA()
			if a != 0 {

				nr, ng, nb := hubble_palette.ChangeHue(r, g, b)

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

func saveAsPng(frame image.Image) error {
	writer, err := os.Create(param.InputParams.OutputDir + "/" + param.InputParams.OutputFileName)

	if err != nil {
		return err
	}

	defer writer.Close()

	var outputImage image.Image
	outputImage = frame

	if param.InputParams.Paletted {
		palettedImg := image.NewPaletted(frame.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(
			palettedImg,
			palettedImg.Bounds(),
			frame,
			frame.Bounds().Min,
		)
		outputImage = palettedImg
	}

	return png.Encode(writer, outputImage)
}

func saveAsWebp(images []image.Image) error {
	writer, err := os.Create(param.InputParams.OutputDir + "/" + param.InputParams.OutputFileName)

	if err != nil {
		return err
	}

	defer writer.Close()

	durations := make([]uint, len(images))
	disposals := make([]uint, len(images))
	for i := range images {
		durations[i] = 100
		disposals[i] = 0 // 0 = keep frame (no disposal)
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
