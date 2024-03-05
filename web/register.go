package web

import (
	"rerere.org/fedi-games/games"
	tictactoe "rerere.org/fedi-games/games/tic-tac-toe"
)

var GameMap = map[string]games.Game{
	tictactoe.TicTacToe{}.Name(): tictactoe.NewTicTacToeGame(),
}
