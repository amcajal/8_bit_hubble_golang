package rngm

import (
	"math/rand"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	KeySmallSpriteMaxLayers   = "SmallSpriteMaxLayers"
	KeyMediumSpriteMaxLayers  = "MediumSpriteMaxLayers"
	KeyLargeSpriteMaxLayers   = "LargeSpriteMaxLayers"
	KeySpecialSpriteMaxLayers = "SpecialSpriteMaxLayers"

	KeySmallSpriteMaxSprites   = "SmallSpriteMaxSprites"
	KeyMediumSpriteMaxSprites  = "MediumSpriteMaxSprites"
	KeyLargeSpriteMaxSprites   = "LargeSpriteMaxSprites"
	KeySpecialSpriteMaxSprites = "SpecialSpriteMaxSprites"

	KeyDefaultChance = "DefaultChance"
	KeySpecialChance = "SpecialChance"
)

// rngConfig holds the active configuration values.
// It is initialised with the default values and can be overridden via LoadConfig.
var rngConfig = map[string]int{
	KeySmallSpriteMaxLayers:   10,
	KeyMediumSpriteMaxLayers:  5,
	KeyLargeSpriteMaxLayers:   2,
	KeySpecialSpriteMaxLayers: 3,

	KeySmallSpriteMaxSprites:   50,
	KeyMediumSpriteMaxSprites:  25,
	KeyLargeSpriteMaxSprites:   10,
	KeySpecialSpriteMaxSprites: 1,

	KeyDefaultChance: 100,
	KeySpecialChance: 50,
}

// Config returns the value for the given key from rngConfig.
func Config(key string) int {
	return rngConfig[key]
}

// LoadConfig reads a TOML config file and overrides the default rngConfig values.
// Keys present in the file replace their default counterparts; missing keys keep
// their defaults.
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var overrides map[string]int
	if _, err := toml.Decode(string(data), &overrides); err != nil {
		return err
	}

	for k, v := range overrides {
		if _, known := rngConfig[k]; known {
			rngConfig[k] = v
		}
	}

	return nil
}

var r *rand.Rand

func Init(seed int64) {
	r = rand.New(rand.NewSource(seed))
}

func Intn(n int) int {
	return r.Intn(n)
}
