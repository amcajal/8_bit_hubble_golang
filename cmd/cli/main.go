package main

import (
	"log"

	"github.com/amcajal/8_bit_hubble_golang/internal/galaxy"
	"github.com/amcajal/8_bit_hubble_golang/internal/param"
	"github.com/amcajal/8_bit_hubble_golang/internal/rngm"
)

func main() {

	// Check parameters are correct
	if err := param.CheckParams(); err != nil {
		log.Fatal(err)
	}

	// Initialize seeded random source (enables deterministic output via -s flag)
	rngm.LoadSeed(param.InputParams.Seed)

	// Generate galaxy
	if err := galaxy.GenerateGalaxy(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Galaxy Generated! Seed: %d", param.InputParams.Seed)
}
