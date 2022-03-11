// Package param (eters) check the user input parameters are correct
package param

import (
	"errors"
	"flag"
	"os"
	"strings"
)

const pngExt string = ".png"

var outputDir string
var pngName string
var seed int64

func CheckParams() error {

	flag.StringVar(&outputDir, "o", "./", "Output directory to save the png")
	flag.StringVar(&pngName, "n", "8bh_galaxy.png", "Name of the png image")
	flag.Int64Var(&seed, "s", 42, "Seed to be used in the image generation")

	flag.Parse()

	if paramError := checkOutputDir(); paramError != nil {
		return paramError
	}

	appendExtension(&pngName)

	return nil
}

// Output directory must exist and have write perms
func checkOutputDir() error {

	status, err := os.Stat(outputDir)
	if err != nil {
		return err
	}

	if !status.IsDir() {
		return errors.New("Provided output dir (" + outputDir + ") is NOT a directory")
	}

	// Probably there is a better way to do this:
	// Isolate owner's read and write bits, and check both are high
	if rw := ((status.Mode()) >> 0x07) & 0x03; rw != 0x03 {
		return errors.New("Provided output dir (" + outputDir + ") DOES NOT HAVE read and write perms")
	}

	return nil
}

func appendExtension(pngName *string) {
	if !strings.HasSuffix(*pngName, pngExt) {
		*pngName = *pngName + pngExt
	}
}
