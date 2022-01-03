// imageWriter replaces image_writer module in C
package imageWriter

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var img *image.NRGBA

func SetDimensions(w int, h int) (rc uint8) {
	img = image.NewNRGBA(image.Rect(0, 0, w, h))
	return 0
}

func SetPixel(x int, y int, hexCode int) (rc uint8) {
	img.Set(x, y, color.NRGBA{
		R: uint8(hexCode >> 16 & 0xFF),
		G: uint8(hexCode >> 8 & 0xFF),
		B: uint8(hexCode & 0xFF),
		A: 255,
	})
	return 0
}

func WriteImage() (rc uint8) {
	f, err := os.Create("image.png")
	
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(f, img)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return 0
}
