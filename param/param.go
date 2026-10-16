// Package param (eters) check the user input parameters are correct
package param

import (
	"errors"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/amcajal/8_bit_hubble_golang/rngm"
)

const pngExt string = ".png"
const webpExt string = ".webp"

type InputParameters struct {
	OutputDir      string
	OutputFileName string
	Seed           int64
	ConfigFile     string
	Animated       bool
	Paletted       bool
}

var InputParams = InputParameters{}

func CheckParams() error {
	flag.StringVar(&InputParams.OutputDir, "o", "./", "Output directory to save the png")
	flag.StringVar(&InputParams.OutputFileName, "n", "8bh_galaxy", "Name of the output image")
	flag.Int64Var(&InputParams.Seed, "s", time.Now().Unix(), "Seed to be used in the image generation")
	flag.StringVar(&InputParams.ConfigFile, "c", "", "Path to a TOML config file for RNG parameters")
	flag.BoolVar(&InputParams.Animated, "a", false, "Generate a WebP animated file (think about it as a \"modern\" GIF)")
	flag.BoolVar(&InputParams.Paletted, "p", false, "Force the output image to show only 256 colors from Plan9 palette")

	flag.Parse()

	if paramError := checkOutputDir(); paramError != nil {
		return paramError
	}

	if paramError := checkConfigFile(); paramError != nil {
		return paramError
	}

	appendExtension()

	return nil
}

func checkConfigFile() error {
	if InputParams.ConfigFile != "" {
		if err := rngm.LoadConfig(InputParams.ConfigFile); err != nil {
			return errors.New("Failed to load config file " + InputParams.ConfigFile + ": " + err.Error())
		}
	}

	return nil
}

// Output directory must exist and have write perms
func checkOutputDir() error {

	status, err := os.Stat(InputParams.OutputDir)
	if err != nil {
		return err
	}

	if !status.IsDir() {
		return errors.New("Provided output dir (" + InputParams.OutputDir + ") is NOT a directory")
	}

	// @TODO this only works if the user launching the app is the owner of the target directory
	// Probably there is a better way to do this:
	// Isolate owner's read and write bits, and check both are high
	if rw := ((status.Mode()) >> 0x07) & 0x03; rw != 0x03 {
		return errors.New("Provided output dir (" + InputParams.OutputDir + ") DOES NOT HAVE read and write perms")
	}

	return nil
}

func appendExtension() {
	if InputParams.Animated {
		if !strings.HasSuffix(InputParams.OutputFileName, webpExt) {
			InputParams.OutputFileName = InputParams.OutputFileName + webpExt
		}
	} else {
		if !strings.HasSuffix(InputParams.OutputFileName, pngExt) {
			InputParams.OutputFileName = InputParams.OutputFileName + pngExt
		}
	}
}
