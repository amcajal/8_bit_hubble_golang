package sprites

import (
	"image"
	"image/png"

	"github.com/amcajal/8_bit_hubble_golang/rngm"
	"github.com/amcajal/pixelart/scale"
)

// Name in honor of Bob Ross
func makeBigDecision(forced, limit int) int {
	if forced != -1 {
		return forced
	}
	return rngm.Intn(limit + 1)
}

func GetNoOfLayers(spriteSize Size) int {

	forcedValue := 0
	defaultLimit := 0

	switch spriteSize {
	case Small:
		forcedValue = rngm.Cfg.NoOfSmallSpriteLayers
		defaultLimit = rngm.Cfg.SmallSpriteMaxLayers
	case Medium:
		forcedValue = rngm.Cfg.NoOfMediumSpriteLayers
		defaultLimit = rngm.Cfg.MediumSpriteMaxLayers
	case Large:
		forcedValue = rngm.Cfg.NoOfLargeSpriteLayers
		defaultLimit = rngm.Cfg.LargeSpriteMaxLayers
	default: // Special
		forcedValue = rngm.Cfg.NoOfSpecialSpriteLayers
		defaultLimit = rngm.Cfg.SpecialSpriteMaxLayers
	}

	return makeBigDecision(forcedValue, defaultLimit)
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

	return makeBigDecision(forcedValue, defaultLimit)
}

func ShouldPaintCosmicLatteBackground() bool {
	chance := rngm.Cfg.CosmicLatteBackgroundChance
	return rngm.Intn(100) <= chance
}

func ShouldPaintNebula() bool {
	chance := rngm.Cfg.NebulaChance
	return rngm.Intn(100) <= chance
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

	img = loadPNG(pool[rngm.Intn(len(pool))].Filename)

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

// loadPNG reads a PNG file from the embedded assets FS and decodes it.
func loadPNG(filename string) image.Image {
	f, err := assetsFS.Open(filename)
	if err != nil {
		panic("sprites: cannot open embedded asset " + filename + ": " + err.Error())
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic("sprites: cannot decode PNG " + filename + ": " + err.Error())
	}
	return img
}
