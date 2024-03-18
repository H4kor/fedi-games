package main

import "rerere.org/fedi-games/games/bunkers"

func main() {
	state := bunkers.NewBunkersGameState()
	for i := 0; i < 2; i++ {
		state.Shots = append(state.Shots, bunkers.Shot{
			StartX: state.PosA,
			StartY: state.Terrain().Height[state.PosA] + 15,
			Angle:  30,
			Vel:    70,
		})

	}

	bunkers.Render(state)

}
