package web

import (
	"github.com/H4kor/fedi-games/games"
	"github.com/H4kor/fedi-games/games/bunkers"
	"github.com/H4kor/fedi-games/games/rps"
	tictactoe "github.com/H4kor/fedi-games/games/tic-tac-toe"
	"github.com/H4kor/fedi-games/internal"
)

func DefaultServer() FediGamesServer {
	gamesList := []games.Game{
		tictactoe.NewTicTacToeGame(),
		rps.NewRockPaperScissorGame(),
		bunkers.NewBunkerGame(),
	}

	engine := internal.NewGameEngine()
	server := NewFediGamesServer(engine, gamesList)
	return server
}
