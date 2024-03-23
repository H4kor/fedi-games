package bunkers

import (
	"bytes"
	"image"
	"image/gif"
	"image/png"

	"github.com/H4kor/fedi-games/internal"
)

// renderStep draws the map on the given step
// step 0 = before any shot is fired
// step 1 = after first shot, first shot is shown
// ...
func renderStep(state BunkersGameState, step int) (*image.Paletted, error) {
	canvas := image.NewPaletted(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: WIDTH, Y: HEIGHT},
		},
		PALETTE,
	)

	state.TerrainAtShot(step).Draw(canvas)
	DrawBunker(state.TerrainAtShot(step), state.PosA, uint(PLAYER_A), canvas)
	DrawBunker(state.TerrainAtShot(step), state.PosB, uint(PLAYER_B), canvas)

	if step != 0 {
		state.Shots[step-1].Draw(state, canvas, step)
	}

	return canvas, nil
}

func RenderAnimation(state BunkersGameState) (string, error) {
	images := make([]*image.Paletted, 0)
	var delays []int
	for i := 0; i < len(state.Shots); i++ {
		img, err := renderStep(state, i)
		if err != nil {
			return "", err
		}
		images = append(images, img)
		delays = append(delays, 200)
	}

	if len(delays) > 0 {
		delays[len(delays)-1] = 1000
	}

	var buffer bytes.Buffer
	err := gif.EncodeAll(
		&buffer,
		&gif.GIF{
			Image: images,
			Delay: delays,
		},
	)
	if err != nil {
		return "", err
	}

	return internal.StoreMedia(buffer.Bytes(), "gif")
}

// Render the state into a png file and save it as as media files
// returns the full url of the file
func Render(state BunkersGameState) (string, error) {
	canvas, err := renderStep(state, len(state.Shots))
	if err != nil {
		return "", err
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
