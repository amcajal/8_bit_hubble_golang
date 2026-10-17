package sprites

import (
	"image"
	"image/color"
	"image/png"

	hubble_palette "github.com/amcajal/8_bit_hubble_golang/palette"
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
func GetSprite(spriteSize Size) image.Image {
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
	entry := pool[index]
	img := loadPNG(entry.Filename)

	if spriteSize == Special {
		img = rescaleSprite(img)
	}

	img = ColorizeSprite(img, rngm.Intn(361))

	if param.InputParams.Animated {
		img = BlinkSprite(img, entry)
	}

	return img
}

// ColorizeSprite rotates the hue of every non-transparent pixel.
func ColorizeSprite(img image.Image, hueDegrees int) image.Image {
	hubble_palette.SetHueRotation(hueDegrees)
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r32, g32, b32, a32 := img.At(x, y).RGBA()
			if a32 == 0 {
				continue
			}
			nr, ng, nb := hubble_palette.ChangeHue(r32, g32, b32)
			a16 := a32 >> 8
			dst.SetNRGBA(x, y, color.NRGBA{
				R: uint8(nr & 0xFF),
				G: uint8(ng & 0xFF),
				B: uint8(nb & 0xFF),
				A: uint8(a16),
			})
		}
	}
	return dst
}

func BlinkSprite(img image.Image, entry Sprite) image.Image {
	if !entry.CanBeAnimated {
		return img
	}
	return ApplyOpacity(img, OpacityValues[entry.OpacityIndex])
}

// ApplyOpacity sets the alpha of every non-transparent pixel to alpha.
func ApplyOpacity(img image.Image, alpha uint8) image.Image {
	// Convert to NRGBA so alpha is stored non-premultiplied.
	// This guarantees that setting A=128 truly means 50% opacity,
	// not a half-brightness dimmed sprite.
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := dst.NRGBAAt(x, y)
			src := img.At(x, y)
			r, g, b, a := src.RGBA() // returns 16-bit premultiplied values
			if a == 0 {
				continue // leave fully-transparent pixels at zero
			}
			// Un-premultiply: scale each channel back to 8-bit full brightness
			a16 := a >> 8 // 0-255, as uint32
			c.R = uint8((r >> 8) * 255 / a16)
			c.G = uint8((g >> 8) * 255 / a16)
			c.B = uint8((b >> 8) * 255 / a16)
			c.A = alpha
			dst.SetNRGBA(x, y, c)
		}
	}
	return dst
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

const (
	MODE_INIT = iota
	MODE_ADVANCE
)

func InitOpacityIndexes() {
	setOpacityIndexes(MODE_INIT)
}

func AdvanceOpacityIndexes() {
	setOpacityIndexes(MODE_ADVANCE)
}

func setOpacityIndexes(mode int) {
	v := len(OpacityValues)
	pools := []*[]Sprite{&smallSprites, &mediumSprites, &bigSprites}
	for _, pool := range pools {
		for i := range *pool {
			switch mode {
			case MODE_INIT:
				(*pool)[i].OpacityIndex = rngm.Intn(v)
			case MODE_ADVANCE:
				(*pool)[i].OpacityIndex = ((*pool)[i].OpacityIndex + 1) % v
			}
		}
	}
}
