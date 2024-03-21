package main

import "github.com/H4kor/fedi-games/games/bunkers"

func main() {
	state := bunkers.NewBunkersGameState("a", "b")
	for i := 0; i < 3; i++ {
		state.Shots = append(state.Shots, bunkers.Shot{
			StartX: state.PosA,
			StartY: state.Terrain().At(state.PosA) + 15,
			Angle:  20 * float64(i),
			Vel:    90,
			Wind:   1,
		})

	}

	bunkers.Render(*state)

}
