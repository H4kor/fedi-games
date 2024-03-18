package bunkers

import (
	"image"
	"image/png"
	"log"
	"os"
)

func Render(state BunkersGameState) {
	canvas := image.NewPaletted(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: WIDTH, Y: HEIGHT},
		},
		PALETTE,
	)

	state.Terrain.Draw(canvas)
	DrawBunker(state.Terrain, state.PosA, uint(PLAYER_A), canvas)
	DrawBunker(state.Terrain, state.PosB, uint(PLAYER_B), canvas)

	if len(state.Shots) != 0 {
		state.Shots[len(state.Shots)-1].Draw(state, canvas)
	}

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, canvas); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func DrawBunker(t Terrain, at int, color uint, canvas *image.Paletted) {
	h := t.Height[at]

	BUNKER_IN_GROUND := -2
	BUNKER_HEIGHT := 15
	BUNKER_WIDTH := 8

	for dx := -BUNKER_WIDTH; dx <= BUNKER_WIDTH; dx++ {
		for dy := BUNKER_IN_GROUND; dy <= BUNKER_HEIGHT; dy++ {
			x := dx + at
			y := dy + h
			// check in bound x
			if x < canvas.Bounds().Min.X || x >= canvas.Bounds().Max.X {
				continue
			}
			// check in bound y
			if y < canvas.Bounds().Min.X || y >= canvas.Bounds().Max.Y {
				continue
			}
			canvas.Set(x, canvas.Rect.Max.Y-y, PALETTE[color])
		}
	}
}
