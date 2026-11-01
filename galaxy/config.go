// Package galaxy is the public API for 8-Bit Hubble Golang.
// It can be imported as a library, called from the CLI, or driven by the GUI.
package galaxy

import (
	"fmt"
	"os"
)

const Version string = "v1.1.0"

// Config holds every parameter that controls galaxy generation.
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
