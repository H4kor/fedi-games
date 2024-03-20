package main

import "rerere.org/fedi-games/games/bunkers"

func main() {
	state := bunkers.NewBunkersGameState("a", "b")
	for i := 0; i < 6; i++ {
		state.Shots = append(state.Shots, bunkers.Shot{
			StartX: state.PosA,
			StartY: state.Terrain().At(state.PosA) + 15,
			Angle:  -30,
			Vel:    90,
			Wind:   5,
		})

	}

	bunkers.Render(*state)

}
