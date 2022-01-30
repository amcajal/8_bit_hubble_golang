package main

import (
	"flag"
	"fmt"
	"8_bit_hubble_golang/user_config_checker"
	"image"
	"image/png"
	"image/draw"
	"math/rand"
	"os"
    "encoding/base64"
    "bytes"
)

//cross_star_10x10_golang.png
var sprite string = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAFXRFWHRUaXRsZQA4Yml0SHViYmxlSW1hZ2UhcnPtAAAAQHRFWHREZXNjcmlwdGlvbgBQcm9qZWN0IFVSTDogaHR0cHM6XFxnaXRodWIuY29tL2FtY2FqYWwvOF9iaXRfaHViYmxliG4p0AAAAPF0RVh0Q29weXJpZ2h0ADgtQml0IEh1YmJsZSBDb3B5cmlnaHQgKEMpIDIwMTggQWxiZXJ0byBNYXJ0aW4gQ2FqYWwKVGhpcyBwcm9ncmFtIGNvbWVzIHdpdGggQUJTT0xVVEVMWSBOTyBXQVJSQU5UWTsKVGhpcyBpcyBmcmVlIHNvZnR3YXJlIGRpc3RyaWJ1dGVkIHVuZGVyIEdOVSBHUEwgdjMuMCBMaWNlbnNlLgpGb3IgbW9yZSBkZXRhaWxzLCB2aXNpdCBodHRwczpcXHd3dy5nbnUub3JnXGxpY2Vuc2VzXGdwbC0zLjAuZW4uaHRtbDrE28YAAAAsdEVYdENyZWF0aW9uIFRpbWUAdmllIDI4IGVuZSAyMDIyIDIzOjEwOjU1ICswMTAwbUFD5AAAAAd0SU1FB+YBHBYLEM2OhcUAAAAJcEhZcwAAHsEAAB7BAcNpVFMAAAAEZ0FNQQAAsY8L/GEFAAAAPUlEQVR42mNgYNjyn4EggKvBpxhN7j8QoApu+Q8RgwBGnDrhwIcRiyDEFEzTGRiYCHuEXKsJeYaU4CEuwAFRtzO3mVJEaQAAAABJRU5ErkJggg=="


func main() {

    // Simulating CLI module from original project
	outputDir := flag.String("o", "./", "Output directory to save the png")
	pngName := flag.String("n", "8bh_galaxy.png", "Name of the png image")
	seed := flag.Int("s", 42, "Seed to be used in the image generation")

	flag.Parse()

	fmt.Printf("Options are: %v %v %v\n", *outputDir, *pngName, *seed)
	
	// Initialize seed
	rand.Seed(42)
	
    user_config_checker.Proto()
    
    // Turn base64 string into png image
    pngSprite := b64ToPng(sprite);
    
    // Create image of fixed dimensions
    dim_x, dim_y := 500, 500
    rgba := image.NewRGBA(image.Rect(0,0,dim_x,dim_y))
    
    // Print 20 sprites in "random" places
    sr := pngSprite.Bounds() // Get the full image as a "source rectangle"
    for i := 0; i < 20; i++ {
        x_c := rand.Intn(dim_x);
        y_c := rand.Intn(dim_y);

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

func b64ToPng(b64string string) image.Image {

    // Do the opposite: turn the string back into the original byte array
    originalData, _ := base64.StdEncoding.DecodeString(b64string)
    
    // Instead of a file descriptor to a file, get a "file descriptor" to the byte slice
    reader := bytes.NewReader(originalData)
    
    // create the in memory structure of the image
    img, _ := png.Decode(reader)
    
    return img
}
