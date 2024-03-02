package tictactoe

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"rerere.org/fedi-games/domain/models"
	"rerere.org/fedi-games/games"
)

func NewTicTacToeGame() games.Game {
	return &TicTacToe{}
}

type TicTacToeState struct {
	Fields  []int // 0 = empty, 1 = PlayerA, 2 = Playerb
	PlayerA string
	PlayerB string
}

type TicTacToe struct {
}

func (*TicTacToe) initState(state *TicTacToeState, msg games.GameMsg) error {
	if len(msg.To) != 1 {
		return errors.New("must mention on other player")
	}
	state.Fields = []int{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}
	state.PlayerA = msg.From
	state.PlayerB = msg.To[0]
	return nil
}

// OnMsg implements games.Game.
func (t *TicTacToe) OnMsg(session *models.GameSession, msg games.GameMsg) (games.GameReply, error) {
	state := session.Data.(TicTacToeState)
	if len(state.Fields) != 9 {
		err := t.initState(&state, msg)
		if err != nil {
			return games.GameReply{}, err
		}
	}

	parts := strings.Split(msg.Msg, " ")
	found := 0
	for _, p := range parts {
		slog.Info("Part", "part", p)
		i, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		if i > 0 && i < 10 {
			found = i
			break
		}
	}
	if found == 0 {
		return games.GameReply{}, errors.New("message must include a field number")
	}
	state.Fields[found-1] = 1

	m := "Field:\n"
	for i, f := range state.Fields {
		if f == 0 {
			m += intToEmoji(i + 1)
		} else if f == 1 {
			m += "❌"
		} else if f == 2 {
			m += "🇴"
		}
		if (i+1)%3 == 0 {
			m += "<br>\n"
		}
	}
	slog.Info("Field message", "msg", m)

	return games.GameReply{
		To:  []string{msg.From},
		Msg: m,
	}, nil
}

func intToEmoji(i int) string {
	switch i {
	case 1:
		return "1️⃣"
	case 2:
		return "2️⃣"
	case 3:
		return "3️⃣"
	case 4:
		return "4️⃣"
	case 5:
		return "5️⃣"
	case 6:
		return "6️⃣"
	case 7:
		return "7️⃣"
	case 8:
		return "8️⃣"
	case 9:
		return "9️⃣"
	default:
		return "#️⃣"
	}
}

// NewState implements games.Game.
func (t *TicTacToe) NewState() interface{} {
	return TicTacToeState{}
}

// Name implements games.Game.
func (TicTacToe) Name() string {
	return "tic-tac-toe"
}
