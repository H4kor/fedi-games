package bunkers

import (
	"image/color"
)

var (
	WIDTH     = 800
	HEIGHT    = 600
	SKY       = 0
	GROUND    = 1
	PLAYER_A  = 2
	PLAYER_B  = 3
	TRAIL     = 4
	EXPLOSION = 5
	PALETTE   = []color.Color{
		color.RGBA{R: 200, G: 200, B: 255, A: 255}, // SKY
		color.RGBA{R: 100, G: 255, B: 100, A: 255}, // GROUND
		color.RGBA{R: 255, G: 0, B: 0, A: 255},     // PLAYER_A
		color.RGBA{R: 0, G: 0, B: 255, A: 255},     // PLAYER_B
		color.RGBA{R: 255, G: 200, B: 0, A: 255},   // TRAIL
		color.RGBA{R: 255, G: 255, B: 255, A: 255}, // EXPLOSION
	}
)
