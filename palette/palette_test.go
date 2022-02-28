package palette

import (
    "fmt"
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"testing"
	"image/color"
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

func TestHueChange(t *testing.T) {
	//cross_star_10x10_golang.png
	var sprite string = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAFXRFWHRUaXRsZQA4Yml0SHViYmxlSW1hZ2UhcnPtAAAAQHRFWHREZXNjcmlwdGlvbgBQcm9qZWN0IFVSTDogaHR0cHM6XFxnaXRodWIuY29tL2FtY2FqYWwvOF9iaXRfaHViYmxliG4p0AAAAPF0RVh0Q29weXJpZ2h0ADgtQml0IEh1YmJsZSBDb3B5cmlnaHQgKEMpIDIwMTggQWxiZXJ0byBNYXJ0aW4gQ2FqYWwKVGhpcyBwcm9ncmFtIGNvbWVzIHdpdGggQUJTT0xVVEVMWSBOTyBXQVJSQU5UWTsKVGhpcyBpcyBmcmVlIHNvZnR3YXJlIGRpc3RyaWJ1dGVkIHVuZGVyIEdOVSBHUEwgdjMuMCBMaWNlbnNlLgpGb3IgbW9yZSBkZXRhaWxzLCB2aXNpdCBodHRwczpcXHd3dy5nbnUub3JnXGxpY2Vuc2VzXGdwbC0zLjAuZW4uaHRtbDrE28YAAAAsdEVYdENyZWF0aW9uIFRpbWUAdmllIDI4IGVuZSAyMDIyIDIzOjEwOjU1ICswMTAwbUFD5AAAAAd0SU1FB+YBHBYLEM2OhcUAAAAJcEhZcwAAHsEAAB7BAcNpVFMAAAAEZ0FNQQAAsY8L/GEFAAAAPUlEQVR42mNgYNjyn4EggKvBpxhN7j8QoApu+Q8RgwBGnDrhwIcRiyDEFEzTGRiYCHuEXKsJeYaU4CEuwAFRtzO3mVJEaQAAAABJRU5ErkJggg=="

    /*
    This functions returns a image.Image, which is of type interface.
    Seems heretic, but an interface can be though as a struct with
    pointers to certain methods, and a "hiden" pointer to the underlying
    type.
    */
	pngSprite := b64ToPng(sprite)
	fmt.Printf("Value of the sprite is %T\n", pngSprite)

	// Create image of fixed dimensions
	dim_x, dim_y := 500, 500
	rgba := image.NewRGBA(image.Rect(0, 0, dim_x, dim_y))

	// Print 20 sprites in "random" places
	sr := pngSprite.Bounds() // Get the full image as a "source rectangle"
	
	// Change the hue of the sprite
	SetHueRotation(250)
	rows, columns := sr.Max.X, sr.Max.Y
	for col := 0; col < columns; col++ {
	    for row :=0; row < rows; row++ {
	        sprite_color := pngSprite.At(row, col)
	        if sprite_color.(color.NRGBA).A != 0 { // CORE CONCEPT: access to the underlying type of the interface, which is color.NRGBA
	            r, g, b, _ := sprite_color.RGBA() // This returns uint32 values
	            nr, ng, nb := ChangeHue(r, g, b)
	            // There must be a clever way to turn uint32 values into a color.NRGBA struct
	            var new_rgb color.NRGBA = color.NRGBA{
	                uint8(nr&0xFF), 
	                uint8(ng&0xFF), 
	                uint8(nb&0xFF), 
	                255}
	            pngSprite.(*image.NRGBA).Set(row, col, new_rgb)
	        }
	    }
	}


	
	for i := 0; i < 20; i++ {
		x_c := rand.Intn(dim_x)
		y_c := rand.Intn(dim_y)

		// Draw the pngSprite
		dp := image.Pt(x_c, y_c) // Point in the destiny image where to start to paint the sprite
		r := image.Rectangle{dp, dp.Add(sr.Size())}
		draw.Draw(rgba, r, pngSprite, sr.Min, draw.Src)
	}

	// Save image
	writer, err := os.Create("./galaxy.png")
	if err != nil {
		fmt.Println("Could not create galaxy.png")
	}
	defer writer.Close()

	png.Encode(writer, rgba)
}
