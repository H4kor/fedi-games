package main

import (
	"rerere.org/fedi-games/games"
	"rerere.org/fedi-games/games/bunkers"
	"rerere.org/fedi-games/games/rps"
	tictactoe "rerere.org/fedi-games/games/tic-tac-toe"
	"rerere.org/fedi-games/internal"
	"rerere.org/fedi-games/web"
)

func main() {
	gamesList := []games.Game{
		tictactoe.NewTicTacToeGame(),
		rps.NewRockPaperScissorGame(),
		bunkers.NewBunkerGame(),
	}

	engine := internal.NewGameEngine()
	server := web.NewFediGamesServer(engine, gamesList)
	err := server.Start()
	if err != nil {
		panic(err)
	}

}
