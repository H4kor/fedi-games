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
