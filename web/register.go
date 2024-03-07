package web

import (
	"rerere.org/fedi-games/games"
	"rerere.org/fedi-games/games/rps"
	tictactoe "rerere.org/fedi-games/games/tic-tac-toe"
)

var GameMap = map[string]games.Game{
	tictactoe.TicTacToe{}.Name():  tictactoe.NewTicTacToeGame(),
	rps.RockPaperScissor{}.Name(): rps.NewRockPaperScissorGame(),
}
