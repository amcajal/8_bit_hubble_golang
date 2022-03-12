package main

import (
	"8_bit_hubble_golang/galaxy"
	"8_bit_hubble_golang/param"
	"log"
	"math/rand"
)

func main() {

	// Check parameters are correct
	if err := param.CheckParams(); err != nil {
		log.Fatal(err)
	}

	// Initialize seed
	rand.Seed(param.Seed)

	// Generate galaxy
	if err := galaxy.GenerateGalaxy(); err != nil {
		log.Fatal(err)
	}

	log.Println("Galaxy Generated!")
}
