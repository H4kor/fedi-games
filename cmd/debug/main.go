package main

import "rerere.org/fedi-games/games/bunkers"

func main() {
	terrain := bunkers.NewTerrain()
	state := bunkers.BunkersGameState{
		Terrain: terrain,
	}

	bunkers.Render(state)

}
