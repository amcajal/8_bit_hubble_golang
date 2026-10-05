package rngm

import (
	"math/rand"
	"os"

	"github.com/BurntSushi/toml"
)

type rngConfig struct {
	// Max number of layers per sprite size
	SmallSpriteMaxLayers   int
	MediumSpriteMaxLayers  int
	LargeSpriteMaxLayers   int
	SpecialSpriteMaxLayers int

	// Max number of sprites (of same size) per layer
	SmallSpriteMaxSprites   int
	MediumSpriteMaxSprites  int
	LargeSpriteMaxSprites   int
	SpecialSpriteMaxSprites int

	// Chances (probability 0-100) of a sprite category being painted at all
	DefaultChance int
	SpecialChance int
}

var Cfg = rngConfig{
	SmallSpriteMaxLayers:   10,
	MediumSpriteMaxLayers:  5,
	LargeSpriteMaxLayers:   2,
	SpecialSpriteMaxLayers: 3,

	SmallSpriteMaxSprites:   50,
	MediumSpriteMaxSprites:  25,
	LargeSpriteMaxSprites:   10,
	SpecialSpriteMaxSprites: 1,

	DefaultChance: 100,
	SpecialChance: 50,
}

type tomlOverrides struct {
	SmallSpriteMaxLayers   *int
	MediumSpriteMaxLayers  *int
	LargeSpriteMaxLayers   *int
	SpecialSpriteMaxLayers *int

	SmallSpriteMaxSprites   *int
	MediumSpriteMaxSprites  *int
	LargeSpriteMaxSprites   *int
	SpecialSpriteMaxSprites *int

	DefaultChance *int
	SpecialChance *int
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var o tomlOverrides
	if _, err := toml.Decode(string(data), &o); err != nil {
		return err
	}

	if o.SmallSpriteMaxLayers != nil {
		Cfg.SmallSpriteMaxLayers = *o.SmallSpriteMaxLayers
	}
	if o.MediumSpriteMaxLayers != nil {
		Cfg.MediumSpriteMaxLayers = *o.MediumSpriteMaxLayers
	}
	if o.LargeSpriteMaxLayers != nil {
		Cfg.LargeSpriteMaxLayers = *o.LargeSpriteMaxLayers
	}
	if o.SpecialSpriteMaxLayers != nil {
		Cfg.SpecialSpriteMaxLayers = *o.SpecialSpriteMaxLayers
	}

	if o.SmallSpriteMaxSprites != nil {
		Cfg.SmallSpriteMaxSprites = *o.SmallSpriteMaxSprites
	}
	if o.MediumSpriteMaxSprites != nil {
		Cfg.MediumSpriteMaxSprites = *o.MediumSpriteMaxSprites
	}
	if o.LargeSpriteMaxSprites != nil {
		Cfg.LargeSpriteMaxSprites = *o.LargeSpriteMaxSprites
	}
	if o.SpecialSpriteMaxSprites != nil {
		Cfg.SpecialSpriteMaxSprites = *o.SpecialSpriteMaxSprites
	}

	if o.DefaultChance != nil {
		Cfg.DefaultChance = *o.DefaultChance
	}
	if o.SpecialChance != nil {
		Cfg.SpecialChance = *o.SpecialChance
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
