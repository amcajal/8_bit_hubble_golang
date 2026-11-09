package main

import (
	"log"
	"os"

	"github.com/amcajal/8_bit_hubble_golang/galaxy"
)

func main() {
	cfg, err := galaxy.ConfigFromArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := galaxy.Generate(cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("Galaxy Generated! Seed: %d", cfg.Seed)
}
