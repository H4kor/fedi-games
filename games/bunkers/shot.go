package bunkers

import (
	"image"
	"math"
)

var EXPLOSION_RADIUS = 20
var GRAVITY = 10.0

type Shot struct {
	StartX int
	StartY int
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
		// if out of bound on x or bottom y -> abort
		if int(x) < 0 || int(x) >= WIDTH || int(y) < 0 {
			return trail, false
		}
		// TODO: check collision with bunker
		// check collision with terrain
		if y <= float64(terrain.Height[int(x)]) {
			break
		}

		trail = append(trail, point{int(x), int(y)})

		x += t * dx
		y += t * dy
		dy -= t * GRAVITY
	}

	return trail, true
}

func (s *Shot) Draw(state BunkersGameState, canvas *image.Paletted) {

	trail, hit_terrain := s.getImpact(state.Terrain())

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
			nt.Height[x+dx] = int(math.Min(
				float64(y)-dy,
				float64(nt.Height[x+dx]),
			))
		}
	}

	return nt
}
