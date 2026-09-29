// sketchs is a utility that exports sprites from the sprites package as
// individual PNG files.
//
// By default (no flags), every available sprite is exported.
// Use -n or --name to export a single sprite by name.
//
// Usage:
//
//	go run utils/sketchs/sketchs.go                     # export all sprites
//	go run utils/sketchs/sketchs.go -n cross_star_5x5   # export one sprite
//	go run utils/sketchs/sketchs.go --name planet        # export one sprite
//
// Output files are named: <sprite_name>_<unix_timestamp>.png
package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/amcajal/8_bit_hubble_golang/sprites"
)

func main() {

	// Parse the optional -n / --name flag from the command-line arguments.
	// We do this by hand to keep things simple and avoid importing "flag".
	spriteName := parseNameFlag(os.Args[1:])

	// A single timestamp shared by all files in this run, so they group nicely.
	timestamp := time.Now().Unix()

	if spriteName != "" {
		exportOneSprite(spriteName, timestamp)
	} else {
		exportAllSprites(timestamp)
	}
}

func parseNameFlag(args []string) string {
	for i, arg := range args {
		if (arg == "-n" || arg == "--name") && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func exportOneSprite(name string, timestamp int64) {
	img, found := sprites.FindSpriteByName(name)
	if !found {
		fmt.Fprintf(os.Stderr, "Error: no sprite named %q exists.\n", name)
		os.Exit(1)
	}

	writePNG(name, img, timestamp)
	fmt.Println("Done!")
}

func exportAllSprites(timestamp int64) {
	allSprites := sprites.GetAllSprites()
	fmt.Printf("Found %d sprites. Writing PNG files...\n", len(allSprites))

	for name, img := range allSprites {
		writePNG(name, img, timestamp)
	}

	fmt.Println("Done!")
}

func writePNG(name string, img image.Image, timestamp int64) {
	filename := fmt.Sprintf("%s_%d.png", name, timestamp)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file %s: %v", filename, err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Fatalf("Could not encode sprite %q as PNG: %v", name, err)
	}

	fmt.Printf("  Written: %s\n", filename)
}
