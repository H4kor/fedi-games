package bunkers

import "math/rand"

type BunkersGame struct{}

type BunkersGameState struct {
	Terrain Terrain
	PosA    int
	PosB    int
	Shots   []Shot
}

func NewBunkersGameState() BunkersGameState {
	return BunkersGameState{
		Terrain: NewTerrain(),
		PosA:    50 + rand.Intn(150),
		PosB:    600 + 50 + rand.Intn(150),
		Shots:   make([]Shot, 0),
	}
}
