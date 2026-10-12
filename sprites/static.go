package sprites

import "embed"

//go:embed assets
var assetsFS embed.FS

var smallSprites = []Sprite{
	{Name: "cross_star_5x5", Filename: "assets/cross_star_5x5.png"},
	{Name: "square_star_2x2", Filename: "assets/square_star_2x2.png"},
	{Name: "square_star_4x4", Filename: "assets/square_star_4x4.png"},
	{Name: "square_star_6x6", Filename: "assets/square_star_6x6.png"},
}

var mediumSprites = []Sprite{
	{Name: "cross_star_10x10", Filename: "assets/cross_star_10x10.png"},
	{Name: "cross_star_12x12", Filename: "assets/cross_star_12x12.png"},
	{Name: "cross_star_16x16", Filename: "assets/cross_star_16x16.png"},
	{Name: "cross_star_20x20", Filename: "assets/cross_star_20x20.png"},
}

var bigSprites = []Sprite{
	{Name: "combined_star_48x48", Filename: "assets/combined_star_48x48.png"},
	{Name: "cross_star_30x30", Filename: "assets/cross_star_30x30.png"},
}

var specialSprites = []Sprite{
	{Name: "planet_50x50", Filename: "assets/planet_50x50.png"},
	{Name: "planet_100x100", Filename: "assets/planet_100x100.png"},
	{Name: "planet_125x125", Filename: "assets/planet_125x125.png"},
}
