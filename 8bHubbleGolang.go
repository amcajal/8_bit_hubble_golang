package main

import (
	"log"

	"github.com/amcajal/8_bit_hubble_golang/galaxy"
	"github.com/amcajal/8_bit_hubble_golang/param"
	"github.com/amcajal/8_bit_hubble_golang/rngm"
)

func main() {

	// Check parameters are correct
	if err := param.CheckParams(); err != nil {
		log.Fatal(err)
	}

	// Initialize seeded random source (enables deterministic output via -s flag)
	rngm.Init(param.Seed)

	// Generate galaxy
	if err := galaxy.GenerateGalaxy(); err != nil {
		log.Fatal(err)
	}

	log.Println("Galaxy Generated!")
}
