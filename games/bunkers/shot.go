package bunkers

import (
	"image"
	"math"
)

type Shot struct {
	StartX int
	StartY int
	Wind   int
	Vel    float64
	Angle  float64
}

type point struct {
	X int
	Y int
}

// getImpact returns trail of a shot and whether the shot landed in terrain
// second return (bool) is HIT TERRAIN
func (s *Shot) getImpact(terrain Terrain) ([]point, bool) {
	x := float64(s.StartX)
	y := float64(s.StartY)
	dx := math.Sin(float64(s.Angle)*math.Pi/180.0) * float64(s.Vel)
	dy := math.Cos(float64(s.Angle)*math.Pi/180.0) * float64(s.Vel)

	t := 1.0 / float64(s.Vel)

	trail := make([]point, 0)

	for {
		// if out of bound bottom y -> abort
		if int(y) < 0 {
			return trail, false
		}
		// if in negative and wind also negative
		if x < 0 && s.Wind <= 0 {
			return trail, false
		}
		// if oob on right and wind positive
		if int(x) >= WIDTH && s.Wind >= 0 {
			return trail, false
		}

		// check collision with terrain
		if y <= float64(terrain.At(int(x))) {
			break
		}

		trail = append(trail, point{int(x), int(y)})

		x += t * dx
		y += t * dy
		dy -= t * GRAVITY
		dx += t * float64(s.Wind)
	}

	return trail, true
}

func (s *Shot) Draw(state BunkersGameState, canvas *image.Paletted, step int) {

	trail, hit_terrain := s.getImpact(state.TerrainAtShot(step - 1))

	for _, p := range trail {
		canvas.Set(p.X, canvas.Rect.Max.Y-p.Y, PALETTE[TRAIL])

	}

	if hit_terrain {

		p := trail[len(trail)-1]
		x := float64(p.X)
		y := float64(p.Y)

		// draw explosion
		for dx := -EXPLOSION_RADIUS; dx <= EXPLOSION_RADIUS; dx++ {
			for dy := -EXPLOSION_RADIUS; dy <= EXPLOSION_RADIUS; dy++ {
				d := math.Sqrt(float64(dx*dx + dy*dy))
				if d <= float64(EXPLOSION_RADIUS) {
					canvas.Set(int(x+float64(dx)), canvas.Rect.Max.Y-int(y+float64(dy)), PALETTE[EXPLOSION])
				}
			}
		}
	}

}

// DestroyTerrain returns a new terrain with the effect of the shot
func (s *Shot) DestroyTerrain(t Terrain) Terrain {
	nt := Terrain{
		Height: t.Height,
	}

	trail, hit_terrain := s.getImpact(t)
	if hit_terrain {
		p := trail[len(trail)-1]
		x := p.X
		y := p.Y
		for dx := -EXPLOSION_RADIUS; dx <= EXPLOSION_RADIUS; dx++ {
			// x² + y² = r² , solve for min y
			// y = sqrt(r² - x²)
			dy := math.Sqrt(float64(EXPLOSION_RADIUS*EXPLOSION_RADIUS - dx*dx))
			nt.Set(x+dx, int(math.Min(
				float64(y)-dy,
				float64(nt.At(x+dx)),
			)))
		}
	}

	return nt
}
