package sprites

import (
	"image"
	"image/png"

	"github.com/amcajal/8_bit_hubble_golang/param"
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

	index := rngm.Intn(len(pool))
	img = loadPNG(pool[index].Filename)

	if spriteSize == Special {
		img = rescaleSprite(img)
	}

	if param.InputParams.Animated && pool[index].CanBeAnimated {
		opacity := opacityValues[pool[index].OpacityIndex]
		img = modifyOpacity(img, opacity)
	}
	return
}

func modifyOpacity(img image.Image, alpha uint8) image.Image {
	bounds := img.Bounds()
	switch v := img.(type) {
	case *image.NRGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := v.NRGBAAt(x, y)
				c.A = alpha
				v.SetNRGBA(x, y, c)
			}
		}
		return v
	case *image.RGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := v.RGBAAt(x, y)
				c.A = alpha
				v.SetRGBA(x, y, c)
			}
		}
		return v
	default:
		panic("modifyOpacity: unsupported image type")
	}
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

// TODO: Merge both methods so they do the op depending on a flag, and checking if canBeAnimated is true
func SetOpacityIndexes() {
	v := len(opacityValues)
	for _, sprite := range smallSprites {
		sprite.OpacityIndex = rngm.Intn(v)
	}
	for _, sprite := range mediumSprites {
		sprite.OpacityIndex = rngm.Intn(v)
	}
	for _, sprite := range bigSprites {
		sprite.OpacityIndex = rngm.Intn(v)
	}
}

func AdvanceOpacityIndexes() {
	for _, sprite := range smallSprites {
		sprite.OpacityIndex = (sprite.OpacityIndex + 1) % len(opacityValues)
	}
	for _, sprite := range mediumSprites {
		sprite.OpacityIndex = (sprite.OpacityIndex + 1) % len(opacityValues)
	}
	for _, sprite := range bigSprites {
		sprite.OpacityIndex = (sprite.OpacityIndex + 1) % len(opacityValues)
	}
}
