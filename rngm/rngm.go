package rngm

import (
	"math/rand"
)

// Max number of layers per sprite size
const SmallSpriteMaxLayers int = 10
const MediumSpriteMaxLayers int = 5
const LargeSpriteMaxLayers int = 2
const SpecialSpriteMaxLayers int = 3

// Max number of sprites (of same size) per layer
const SmallSpriteMaxSprites int = 50
const MediumSpriteMaxSprites int = 25
const LargeSpriteMaxSprites int = 10
const SpecialSpriteMaxSprites int = 1

// Chances (or probability) of a sprite category being painted at all
const DefaultChance int = 100
const SpecialChance int = 50

var r *rand.Rand

func Init(seed int64) {
	r = rand.New(rand.NewSource(seed))

	// TODO create map with the contstants above, keeping the names and values
	// name of the map rngConfig
}

func Intn(n int) int {
	return r.Intn(n)
}
