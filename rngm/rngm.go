// Package rngm (RNG Manipulator) holds the parameters that control how many
// sprites are generated per layer, how many layers exist, and the probability
// of each sprite category appearing. These values act as the defaults when no
// .config file is provided by the user.
package rngm

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
const DefChance int = 100
const SpecialChance int = 50
