package sprites

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"github.com/amcajal/8_bit_hubble_golang/rngm"
	"github.com/amcajal/pixelart/scale"
)

func GetNoOfLayers(spriteSize Size) int {

	forcedValue := 0
	defaultLimit := 0

	switch spriteSize {
	case Small:
		forcedValue = rngm.Cfg.NoOfSmallSpriteLayers
		defaultLimit = rngm.Cfg.SmallSpriteMaxLayers
	case Medium:
		forcedValue = rngm.Cfg.NoOfSmallSpriteLayers
		defaultLimit = rngm.Cfg.MediumSpriteMaxLayers
	case Large:
		forcedValue = rngm.Cfg.LargeSpriteMaxLayers
		defaultLimit = rngm.Cfg.LargeSpriteMaxLayers
	default: // Special
		forcedValue = rngm.Cfg.SpecialSpriteMaxLayers
		defaultLimit = rngm.Cfg.SpecialSpriteMaxLayers
	}

	if forcedValue != -1 {
		return forcedValue
	} else {
		return rngm.Intn(defaultLimit + 1)
	}
}

func GetNoOfElements(spriteSize Size) int {

	forcedValue := 0
	defaultLimit := 0

	switch spriteSize {
	case Small:
		forcedValue = rngm.Cfg.NoOfSmallSpriteElements
		defaultLimit = rngm.Cfg.SmallSpriteMaxSprites
	case Medium:
		forcedValue = rngm.Cfg.NoOfMediumSpriteElements
		defaultLimit = rngm.Cfg.MediumSpriteMaxSprites
	case Large:
		forcedValue = rngm.Cfg.NoOfLargeSpriteElements
		defaultLimit = rngm.Cfg.LargeSpriteMaxSprites
	default: // Special
		forcedValue = rngm.Cfg.NoOfSpecialSpriteElements
		defaultLimit = rngm.Cfg.SpecialSpriteMaxSprites
	}

	if forcedValue != -1 {
		return forcedValue
	} else {
		return rngm.Intn(defaultLimit + 1)
	}
}

func ShouldPaintSprite(spriteSize Size) bool {
	chance := 0

	switch spriteSize {
	case Special:
		chance = rngm.Cfg.SpecialChance
	default:
		chance = rngm.Cfg.DefaultChance
	}
	return rngm.Intn(100) <= chance
}

func AllSprites() []Sprite {
	all := make([]Sprite, 0, len(smallSprites)+len(mediumSprites)+len(bigSprites)+len(specialSprites))
	all = append(all, smallSprites...)
	all = append(all, mediumSprites...)
	all = append(all, bigSprites...)
	all = append(all, specialSprites...)
	return all
}

func GetSprite(spriteSize Size) (img image.Image) {
	var pool []Sprite

	switch spriteSize {
	case Small:
		pool = smallSprites
	case Medium:
		pool = mediumSprites
	case Large:
		pool = bigSprites
	default: // Special
		pool = specialSprites
	}

	img = base64ToPng(pool[rngm.Intn(len(pool))].B64Value)

	if spriteSize == Special {
		img = rescaleSprite(img)
	}

	return
}

func rescaleSprite(sprite image.Image) (newSprite image.Image) {
	newSprite = sprite

	if p := rngm.Intn(4); p < 2 {
		for i := 0; i <= p; i++ {
			newSprite = scale.Scale2X(sprite.(*image.NRGBA))
		}
	}
	return
}

func base64ToPng(b64string string) image.Image {

	// Turn the string back into the original byte array
	originalData, _ := base64.StdEncoding.DecodeString(b64string)

	// Instead of a file descriptor to a file, get a "file descriptor" to the byte slice
	reader := bytes.NewReader(originalData)

	// Create the in-memory image structure
	img, _ := png.Decode(reader)

	return img
}
