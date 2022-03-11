// Package param (eters) check the user input parameters are correct
package param

import (
	"errors"
	"flag"
	"os"
	"strings"
)

const pngExt string = ".png"

var OutputDir string
var PngName string
var Seed int64

func CheckParams() error {

	flag.StringVar(&OutputDir, "o", "./", "Output directory to save the png")
	flag.StringVar(&PngName, "n", "8bh_galaxy.png", "Name of the png image")
	flag.Int64Var(&Seed, "s", 42, "Seed to be used in the image generation")

	flag.Parse()

	if paramError := checkOutputDir(); paramError != nil {
		return paramError
	}

	appendExtension()

	return nil
}

// Output directory must exist and have write perms
func checkOutputDir() error {

	status, err := os.Stat(OutputDir)
	if err != nil {
		return err
	}

	if !status.IsDir() {
		return errors.New("Provided output dir (" + OutputDir + ") is NOT a directory")
	}

	// @TODO this only works if the user launching the app is the owner of the target directory
	// Probably there is a better way to do this:
	// Isolate owner's read and write bits, and check both are high
	if rw := ((status.Mode()) >> 0x07) & 0x03; rw != 0x03 {
		return errors.New("Provided output dir (" + OutputDir + ") DOES NOT HAVE read and write perms")
	}

	return nil
}

func appendExtension() {
	if !strings.HasSuffix(PngName, pngExt) {
		PngName = PngName + pngExt
	}
}
