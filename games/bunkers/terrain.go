package bunkers

import (
	"image"
	"math"
	"math/rand"
)

type Terrain struct {
	Height []int
}

func (t *Terrain) Draw(img *image.Paletted) {

	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		h := t.Height[x]
		for y := img.Rect.Max.Y - 1; y >= img.Rect.Min.Y; y-- {
			if HEIGHT-y > h {
				break
			}
			img.Set(x, y, PALETTE[GROUND])
		}
	}
}

func NewTerrain() Terrain {
	heights := make([]int, WIDTH)

	WAVES := 6
	amps := make([]float64, WAVES)
	freqs := make([]float64, WAVES)
	offs := make([]float64, WAVES)

	rem_h := float64(HEIGHT) / 4

	for i := 0; i < WAVES; i++ {
		amps[i] = rand.Float64() * rem_h / float64(WAVES) * 2
		rem_h -= amps[i]
		freqs[i] = rand.Float64() * 2 * math.Pi / (float64(WIDTH) * 0.25)
		offs[i] = rand.Float64() * float64(WIDTH)
	}

	for i := 0; i < WIDTH; i++ {
		heights[i] = HEIGHT / 4
		heights[i] += int(math.Cos(float64(i-WIDTH/2)*math.Pi/float64(WIDTH/2)) * float64(HEIGHT/2))
		heights[i] = int(math.Max(float64(heights[i]), 100))
		for k := 0; k < WAVES; k++ {
			heights[i] += int(amps[k] * math.Sin(freqs[k]*(float64(i)+offs[k])))
		}
		heights[i] = int(math.Max(float64(heights[i]), 50))
	}

	return Terrain{
		Height: heights,
	}
}
