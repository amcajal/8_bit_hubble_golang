package main

import (
	"8_bit_hubble_golang/param"
	"8_bit_hubble_golang/sprites"
	"image"
	"image/png"
	"image/draw"
	"math/rand"
	"os"
    "log"
)


func main() {

    // Check parameters are correct
    err := param.CheckParams()
    if err != nil {
        log.Fatal(err)
    }
    
	// Initialize seed
	rand.Seed(param.Seed)
    
    // Turn base64 string into png image
    pngSprite := sprites.GetSprite(sprites.Small)
    
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
    writer, err := os.Create(param.OutputDir + "/" + param.PngName)
    if err != nil {
        log.Fatal(err)
    }
    defer writer.Close()
    
    png.Encode(writer, rgba)
}
