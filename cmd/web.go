package main

import (
	"github.com/H4kor/fedi-games/games"
	"github.com/H4kor/fedi-games/games/bunkers"
	"github.com/H4kor/fedi-games/games/rps"
	tictactoe "github.com/H4kor/fedi-games/games/tic-tac-toe"
	"github.com/H4kor/fedi-games/internal"
	"github.com/H4kor/fedi-games/web"
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
