// Package galaxy is the public API for 8-Bit Hubble Golang.
// It can be imported as a library, called from the CLI, or driven by the GUI.
package galaxy

import (
	"fmt"
	"os"
)

const Version string = "v2.0.0"

// RenderConfig holds the parameters that control image generation.
// Pass this to ApplyConfig (to seed the RNG) and then to RenderFrames.
//
// Library consumers are responsible for populating this struct correctly and
// for calling ApplyConfig before RenderFrames — see those functions for details.
type RenderConfig struct {
	Width  int
	Height int
	Seed   int64

	// Animated controls whether RenderFrames produces a single frame (false)
	// or the full sequence of frames required for an animated WebP (true).
	Animated bool

	// RNGConfigFile is an optional path to a TOML file that overrides sprite
	// probability parameters. Leave empty to use built-in defaults.
	// Note: RNGConfigFile is loaded by ApplyConfig, not by RenderFrames.
	RNGConfigFile string

	// RetainLayers instructs RenderFrames to populate GalaxyData.Layers with
	// the individual compositing layers for each frame (background, small
	// sprites, medium sprites, large sprites, special sprites).
	//
	// This is intended for library consumers who need per-layer access, e.g.
	// for parallax scrolling in 2D game engines. It is NOT exposed through the
	// CLI or GUI because it can consume significant memory for large or
	// animated outputs.
	RetainLayers bool
}

// OutputConfig holds the parameters that control how a generated image is
// written to disk. Pass this to SaveAsPng or SaveAsWebp.
type OutputConfig struct {
	OutputDir      string
	OutputFileName string

	// Paletted restricts PNG output to 256 colours from the Plan 9 palette
	// using Floyd-Steinberg dithering. Ignored by SaveAsWebp.
	Paletted bool
}

type Config struct {
	RenderConfig
	OutputConfig
}

// Validate checks that the Config is consistent and that referenced paths
// exist and are accessible. It is called automatically by Generate; library
// consumers using RenderFrames directly may call it themselves if desired.
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
