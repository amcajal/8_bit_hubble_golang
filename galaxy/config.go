// Package galaxy is the public API for 8-Bit Hubble Golang.
// It can be imported as a library, called from the CLI, or driven by the GUI.
package galaxy

import (
	"fmt"
	"os"
)

const Version string = "v2.0.0"

// RenderConfig controls image generation. Pass it to ApplyConfig first, then
// to RenderFrames.
type RenderConfig struct {
	Width  int
	Height int
	Seed   int64 // Same seed always produces the same image.

	// Animated: false -> single PNG frame; true -> full WebP animation sequence.
	Animated bool

	// RNGConfigFile is an optional TOML file that overrides sprite probability
	// parameters. Leave empty to use built-in defaults.
	// Loaded by ApplyConfig, not by RenderFrames.
	RNGConfigFile string

	// RetainLayers populates GalaxyData.Layers with the individual compositing
	// layers for each frame, useful for parallax scrolling in game engines.
	// Not available through the CLI or GUI — can consume significant memory.
	RetainLayers bool
}

// OutputConfig controls how generated images are written to disk.
// Pass it to SaveAsPng or SaveAsWebp.
type OutputConfig struct {
	OutputDir      string
	OutputFileName string

	// Paletted restricts PNG output to 256 colours (Plan 9 palette,
	// Floyd-Steinberg dithering). Ignored by SaveAsWebp.
	Paletted bool
}

// Config is the all-in-one configuration used by the CLI and GUI.
// It embeds RenderConfig and OutputConfig. Library consumers should prefer
// working with the two sub-structs directly.
type Config struct {
	RenderConfig
	OutputConfig
}

// Validate checks that the Config fields are consistent and that all
// referenced paths exist and are accessible. Called automatically by Generate.
func (c Config) Validate() error {
	if c.Width <= 0 || c.Height <= 0 {
		return fmt.Errorf("width and height must be positive (got %dx%d)", c.Width, c.Height)
	}

	if c.OutputFileName == "" {
		return fmt.Errorf("output file name must not be empty")
	}

	if c.RNGConfigFile != "" {
		if _, err := os.Stat(c.RNGConfigFile); err != nil {
			return fmt.Errorf("RNG config file %q not found: %w", c.RNGConfigFile, err)
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
