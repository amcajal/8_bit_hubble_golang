package main

import (
	"flag"
	"fmt"
)

func main() {


    // Simulating CLI module from original project
	outputDir := flag.String("o", "./", "Output directory to save the png")
	pngName := flag.String("n", "8bh_galaxy.png", "Name of the png image")
	seed := flag.Int("s", 42, "Seed to be used in the image generation")

	flag.Parse()

	fmt.Printf("Options are: %v %v %v\n", *outputDir, *pngName, *seed)
	
	
}
