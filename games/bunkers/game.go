package bunkers

import "math/rand"

type BunkersGame struct{}

type BunkersGameState struct {
	InitTerrain Terrain
	PosA        int
	PosB        int
	Shots       []Shot
}

func NewBunkersGameState() BunkersGameState {
	return BunkersGameState{
		InitTerrain: NewTerrain(),
		PosA:        50 + rand.Intn(150),
		PosB:        600 + 50 + rand.Intn(150),
		Shots:       make([]Shot, 0),
	}
}

func (s *BunkersGameState) Terrain() Terrain {
	t := s.InitTerrain.Copy()
	if len(s.Shots) > 1 {
		for _, shot := range s.Shots[:len(s.Shots)-1] {
			t = shot.DestroyTerrain(t)
		}
	}
	return t
}
