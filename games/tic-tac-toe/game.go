package tictactoe

import "rerere.org/fedi-games/games"

func NewTicTacToeGame() games.Game {
	return &TicTacToe{}
}

type TicTacToe struct {
}

// OnMsg implements games.Game.
func (t *TicTacToe) OnMsg(msg games.GameMsg) (games.GameReply, error) {
	m := ""
	m += "1️⃣2️⃣3️⃣<br>\n"
	m += "4️⃣5️⃣6️⃣<br>\n"
	m += "7️⃣8️⃣9️⃣<br>\n"

	return games.GameReply{
		To:  []string{msg.From},
		Msg: m,
	}, nil
}
