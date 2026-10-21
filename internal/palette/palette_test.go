package palette

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func b64ToPng(b64string string) image.Image {

	// Do the opposite: turn the string back into the original byte array
	originalData, _ := base64.StdEncoding.DecodeString(b64string)

	// Instead of a file descriptor to a file, get a "file descriptor" to the byte slice
	reader := bytes.NewReader(originalData)

	// create the in memory structure of the image
	img, _ := png.Decode(reader)

	return img
}

// Surely there is a way to do this thru color.Color interface and the
// "rgb" types in image/color package. But for simplicity, it is performed
// manually
type testColor struct {
	r, g, b, a uint8
}

func compareRGBA(inputColor color.Color, expColor testColor) (testColor, bool) {

	var r, g, b, a uint8

	switch inputColor.(type) {
	case color.NRGBA:
		r, g, b, a = inputColor.(color.NRGBA).R, inputColor.(color.NRGBA).G, inputColor.(color.NRGBA).B, inputColor.(color.NRGBA).A
		//fmt.Println("Is NRGBA")
	case color.RGBA:
		r, g, b, a = inputColor.(color.RGBA).R, inputColor.(color.RGBA).G, inputColor.(color.RGBA).B, inputColor.(color.RGBA).A
		//fmt.Println("Is RGBA")
	default:
		panic("Error: Image color type is not supported")
	}

	// Another way to obtain the values:
	// n, r, g, b, a := color.RGBA()
	// And now, divide each one by 0x101

	tc := testColor{r, g, b, a}
	return tc, ((r == expColor.r) &&
		(g == expColor.g) &&
		(b == expColor.b) &&
		(a == expColor.a))
}

func colorizeSprite(sprite *image.Image) {

	sb := (*sprite).Bounds()

	// Colorize pixels of the sprite
	rows, columns := sb.Max.X, sb.Max.Y
	for col := 0; col < columns; col++ {
		for row := 0; row < rows; row++ {
			sprite_color := (*sprite).At(row, col)
			r, g, b, a := sprite_color.RGBA()
			if a != 0 {

				nr, ng, nb := ChangeHue(r, g, b)

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

func TestHueChange(t *testing.T) {
	/*
	   The next string is the base64 representation of a well known 2x2 PNG
	   image.

	   Consider the 2x2 image as follows:
	   AB
	   CD

	   The RGB values of each pixel (A, B, C and D) are chosen as follows:
	   A: 0,0,255
	   B: 188, 189, 254
	   C: 2, 6, 108
	   D: transparent (Alpha = 0)


	   Given those colors, and performing a manual hue change in several
	   resources (GIMP, online converters, python code from
	   https://stackoverflow.com/questions/8507885/shift-hue-of-an-rgb-color),
	   in theory a HUE change to 250 should give the following outputs:
	   A: 0, 252, 27
	   B: 183, 253, 195
	   C: 0, 107, 13
	   D: transparent (Alpha = 0)

	   This test checks this transformation
	*/
	var b64Sprite string = "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAALHRFWHRDcmVhdGlvbiBUaW1lAGRvbSAxMyBtYXIgMjAyMiAwMDowMjo1OCArMDEwMBQ7QBcAAAAHdElNRQfmAwwXBwRQ+nB8AAAACXBIWXMAAB9AAAAfQAGTqFunAAAABGdBTUEAALGPC/xhBQAAABlJREFUeNpjZGD4/3/P3v8MDExsOf+BgAEASxgI6PchfpcAAAAASUVORK5CYII="

	// Turn sprite into image
	sprite := b64ToPng(b64Sprite)

	t.Log("FIRST ROUND OF CHECKS: ORIGINAL IMAGE\n==========\n")

	// Check values of each pixel
	// pXc stands for "pixel X comparisson"
	// It seems to be so counter-intuitive that x,y coordenates doesnt work
	// as a matrix, Instead, they work as a coordenate char
	expectedColor1 := testColor{0, 0, 255, 255}
	if tc, p1c := compareRGBA(sprite.At(0, 0), expectedColor1); p1c != true {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n in pixel 0,0, expected %v", tc, expectedColor1)
		t.Fail()
	}

	expectedColor2 := testColor{188, 189, 254, 255}
	if tc, p2c := compareRGBA(sprite.At(1, 0), expectedColor2); p2c != true {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n in pixel 0,1, expected %v", tc, expectedColor2)
		t.Fail()
	}

	expectedColor3 := testColor{2, 6, 108, 255}
	if tc, p3c := compareRGBA(sprite.At(0, 1), expectedColor3); p3c != true {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n in pixel 1,0, expected %v", tc, expectedColor3)
		t.Fail()
	}

	if _, _, _, a := sprite.At(1, 1).RGBA(); a != 0 {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n (alpha value) in pixel 1, 1", a)
	}

	// Change hue
	SetHueRotation(250)
	colorizeSprite(&sprite)

	// Check values of each pixel. They should have changed
	t.Log("SECOND ROUND OF CHECKS: COLORIZED IMAGE\n==========\n")

	// Check values of each pixel
	// pXc stands for "pixel X comparisson"
	expectedColor1 = testColor{0, 252, 27, 255}
	if tc, p1c := compareRGBA(sprite.At(0, 0), expectedColor1); p1c != true {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n in pixel 0,0, expected %v", tc, expectedColor1)
		t.Fail()
	}

	expectedColor2 = testColor{183, 253, 195, 255}
	if tc, p2c := compareRGBA(sprite.At(1, 0), expectedColor2); p2c != true {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n in pixel 0,1, expected %v", tc, expectedColor2)
		t.Fail()
	}

	expectedColor3 = testColor{0, 107, 13, 255}
	if tc, p3c := compareRGBA(sprite.At(0, 1), expectedColor3); p3c != true {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n in pixel 1,0, expected %v", tc, expectedColor3)
		t.Fail()
	}

	if _, _, _, a := sprite.At(1, 1).RGBA(); a != 0 {
		t.Logf("FAILURE: Input image has no correct colors. Got %v\n (alpha value) in pixel 1, 1", a)
	}
}
