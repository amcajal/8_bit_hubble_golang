// Package param (eters) check the user input parameters are correct
package param

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amcajal/8_bit_hubble_golang/rngm"
)

const version string = "v1.0.0"
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
	// -h / --help and -v / --version are handled before flag.Parse because
	// the flag package intercepts -h itself and would print its own usage.
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "-v", "--version":
			fmt.Println(version)
			os.Exit(0)
		}
	}

	// Register each flag twice (short + long) sharing the same variable.
	defaultSeed := time.Now().Unix()

	flag.StringVar(&InputParams.OutputDir, "o", "./", "")
	flag.StringVar(&InputParams.OutputDir, "output", "./", "")

	flag.StringVar(&InputParams.OutputFileName, "n", "8bh_galaxy", "")
	flag.StringVar(&InputParams.OutputFileName, "name", "8bh_galaxy", "")

	flag.Int64Var(&InputParams.Seed, "s", defaultSeed, "")
	flag.Int64Var(&InputParams.Seed, "seed", defaultSeed, "")

	flag.StringVar(&InputParams.ConfigFile, "c", "", "")
	flag.StringVar(&InputParams.ConfigFile, "config", "", "")

	flag.BoolVar(&InputParams.Animated, "a", false, "")
	flag.BoolVar(&InputParams.Animated, "animated", false, "")

	flag.BoolVar(&InputParams.Paletted, "p", false, "")
	flag.BoolVar(&InputParams.Paletted, "paletted", false, "")

	// Suppress the default usage output
	flag.Usage = func() {}

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

func printHelp() {
	fmt.Printf(`8-bit Hubble Golang %s

Generates a procedural "8-bit galaxy" image (PNG) or animation (WebP).

USAGE:
    8bHubbleGolang [OPTIONS]

OPTIONS:
    -o, --output    <path>   Output directory where the image will be saved.
                             (default: current directory)

    -n, --name      <name>   Base name of the output file, without extension.
                             The correct extension (.png or .webp) is appended
                             automatically.
                             (default: 8bh_galaxy)

    -s, --seed      <int>    Seed for the random number generator. Using the
                             same seed produces the same image every time.
                             (default: current Unix timestamp)

    -c, --config    <path>   Path to a TOML config file to override RNG
                             parameters such as sprite counts and chances.

    -a, --animated           Generate a WebP animated file instead of a static
                             PNG. Sprites blink by cycling through opacity
                             levels across frames.

    -p, --paletted           Restrict the output image to 256 colours from the
                             Plan 9 palette (Floyd-Steinberg dithering).

    -v, --version            Print the version number and exit.

    -h, --help               Show this help message and exit.

EXAMPLES:
    8bHubbleGolang
    8bHubbleGolang -s 42 -o ./output -n my_galaxy
    8bHubbleGolang --seed 42 --output ./output --name my_galaxy
    8bHubbleGolang -a -s 99 -n blinking_stars
`, version)
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

	tmp, err := os.CreateTemp(InputParams.OutputDir, ".8bh_write_check_*")
	if err != nil {
		return errors.New("Provided output dir (" + InputParams.OutputDir + ") is not writable: " + err.Error())
	}
	tmp.Close()
	os.Remove(tmp.Name())

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
