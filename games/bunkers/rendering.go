package bunkers

import (
	"bytes"
	"image"
	"image/png"

	"github.com/H4kor/fedi-games/internal"
)

// Render the state into a png file and save it as as media files
// returns the full url of the file
func Render(state BunkersGameState) (string, error) {
	canvas := image.NewPaletted(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: WIDTH, Y: HEIGHT},
		},
		PALETTE,
	)

	state.Terrain().Draw(canvas)
	DrawBunker(state.Terrain(), state.PosA, uint(PLAYER_A), canvas)
	DrawBunker(state.Terrain(), state.PosB, uint(PLAYER_B), canvas)

	if len(state.Shots) != 0 {
		state.Shots[len(state.Shots)-1].Draw(state, canvas)
	}

	var buffer bytes.Buffer

	if err := png.Encode(&buffer, canvas); err != nil {
		return "", err
	}

	return internal.StoreMedia(buffer.Bytes(), "png")
}

func DrawBunker(t Terrain, at int, color uint, canvas *image.Paletted) {
	h := t.At(at)

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
