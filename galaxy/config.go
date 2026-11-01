// Package galaxy is the public API for 8-Bit Hubble Golang.
// It can be imported as a library, called from the CLI, or driven by the GUI.
package galaxy

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const Version string = "v1.0.0"

type Config struct {
	Width  int
	Height int
	Seed   int64

	ConfigFile string
	Animated   bool
	Paletted   bool

	OutputDir      string
	OutputFileName string
}

func (c Config) Validate() error {
	if c.Width <= 0 || c.Height <= 0 {
		return fmt.Errorf("width and height must be positive (got %dx%d)", c.Width, c.Height)
	}

	if c.OutputFileName == "" {
		return fmt.Errorf("output file name must not be empty")
	}

	if c.ConfigFile != "" {
		if _, err := os.Stat(c.ConfigFile); err != nil {
			return fmt.Errorf("config file %q not found: %w", c.ConfigFile, err)
		}
	}

	if c.OutputDir != "" {
		info, err := os.Stat(c.OutputDir)
		if err != nil {
			return fmt.Errorf("output directory %q does not exist: %w", c.OutputDir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("output directory %q is not a directory", c.OutputDir)
		}
		tmp, err := os.CreateTemp(c.OutputDir, ".8bh_write_check_*")
		if err != nil {
			return fmt.Errorf("output directory %q is not writable: %w", c.OutputDir, err)
		}
		tmp.Close()
		os.Remove(tmp.Name())
	}

	return nil
}

// ConfigFromArgs parses os.Args into a Config, handling -h/--help and
// -v/--version as early exits. Returns an error if flag parsing fails.
func ConfigFromArgs(args []string) (Config, error) {
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			PrintHelp()
			os.Exit(0)
		case "-v", "--version":
			fmt.Println(Version)
			os.Exit(0)
		}
	}

	// Use a fresh FlagSet so re-entrant calls don't panic with "flag redefined".
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.Usage = func() {}

	cfg := Config{Seed: time.Now().Unix()}

	fs.StringVar(&cfg.OutputDir, "o", "./", "")
	fs.StringVar(&cfg.OutputDir, "output", "./", "")

	fs.StringVar(&cfg.OutputFileName, "n", "8bh_galaxy", "")
	fs.StringVar(&cfg.OutputFileName, "name", "8bh_galaxy", "")

	fs.Int64Var(&cfg.Seed, "s", cfg.Seed, "")
	fs.Int64Var(&cfg.Seed, "seed", cfg.Seed, "")

	fs.StringVar(&cfg.ConfigFile, "c", "", "")
	fs.StringVar(&cfg.ConfigFile, "config", "", "")

	fs.BoolVar(&cfg.Animated, "a", false, "")
	fs.BoolVar(&cfg.Animated, "animated", false, "")

	fs.BoolVar(&cfg.Paletted, "p", false, "")
	fs.BoolVar(&cfg.Paletted, "paletted", false, "")

	fs.IntVar(&cfg.Width, "W", 500, "")
	fs.IntVar(&cfg.Width, "width", 500, "")

	fs.IntVar(&cfg.Height, "H", 500, "")
	fs.IntVar(&cfg.Height, "height", 500, "")

	if err := fs.Parse(args[1:]); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func PrintHelp() {
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

    -W, --width     <int>    Width of the output image in pixels.
                             (default: 500)

    -H, --height    <int>    Height of the output image in pixels.
                             (default: 500)

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
`, Version)
}
