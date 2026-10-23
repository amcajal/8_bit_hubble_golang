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

	"github.com/amcajal/8_bit_hubble_golang/internal/param"
	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
	"github.com/amcajal/8_bit_hubble_golang/internal/sprites"
	"github.com/amcajal/bild_modfor8bhg/noise"
)

// Fixed dimensions of the output image
const dim_x int = 500 // width
const dim_y int = 500 // height

func GenerateGalaxy() error {
	if param.InputParams.Animated {
		return saveAsWebp(GenerateFrames())
	} else {
		return saveAsPng(GenerateSingleImage())
	}
}

func paintForeground(canvas *image.NRGBA) {
	rngm.ReloadSeed() // Reset seed so output in animated and static modes match
	paintSprite(canvas, sprites.Small)
	paintSprite(canvas, sprites.Medium)
	paintSprite(canvas, sprites.Large)
	paintSprite(canvas, sprites.Special)
}

func GenerateSingleImage() image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, dim_x, dim_y))
	noise.SetRandSource(rngm.GetRand())
	paintBackground(canvas, buildBackground())
	paintForeground(canvas)
	return canvas
}

func GenerateFrames() []image.Image {
	noOfFrames := len(sprites.OpacityValues)
	images := make([]image.Image, noOfFrames)

	noise.SetRandSource(rngm.GetRand())
	bg := buildBackground()

	sprites.InitOpacityIndexes()
	for i := 0; i < noOfFrames; i++ {
		frame := image.NewNRGBA(image.Rect(0, 0, dim_x, dim_y))
		paintBackground(frame, bg)
		paintForeground(frame)
		images[i] = frame
		sprites.AdvanceOpacityIndexes()
	}
	return images
}

func buildBackground() image.Image {
	bg := image.NewNRGBA(image.Rect(0, 0, dim_x, dim_y))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(bg, bg.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	if sprites.ShouldPaintNebula() {
		nebula := sprites.GenerateNebula()
		draw.Draw(bg, bg.Bounds(), nebula, image.Point{}, draw.Over)
	}

	if sprites.ShouldPaintCosmicLatteBackground() {
		stars_background := sprites.GetCosmicLatteBackground()
		draw.Draw(bg, bg.Bounds(), stars_background, image.Point{}, draw.Over)
	}

	return bg
}

func paintBackground(frame *image.NRGBA, bg image.Image) {
	draw.Draw(frame, frame.Bounds(), bg, image.Point{}, draw.Src)
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

		img := sprites.GetSprite(spriteSize)

		// sb stands for "sprite bounds"
		sb := img.Bounds()

		for e := 1; e <= elements; e++ {

			// Random position (coordinates) in the image
			x_c := rngm.Intn(dim_x)
			y_c := rngm.Intn(dim_y)

			// Draw the sprite
			dp := image.Pt(x_c, y_c)
			r := image.Rectangle{dp, dp.Add(sb.Size())}
			draw.Draw(frame, r, img, sb.Min, draw.Over)
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
		durations[i] = 175
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
