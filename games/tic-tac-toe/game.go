package tictactoe

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/H4kor/fedi-games/domain/models"
	"github.com/H4kor/fedi-games/games"
	"github.com/H4kor/fedi-games/internal/acpub"
)

func NewTicTacToeGame() games.Game {
	return &TicTacToe{}
}

type TicTacToeState struct {
	Fields   []int // 0 = empty, 1 = PlayerA, 2 = PlayerB
	PlayerA  string
	PlayerB  string
	WhosTurn int // 1 = PlayerA, 2 = PlayerB
	Ended    bool
}

type TicTacToe struct {
}

// Example implements games.Game.
func (t *TicTacToe) Example() string {
	return "@alice@example.com 5"
}

// Summary implements games.Game.
func (t *TicTacToe) Summary() string {
	return `
Classic game of tic-tac-toe. <br>
Mention me and a component to start a game.
Moves are selected by replying with a number between 1 and 9.
`

}

func (*TicTacToe) initState(state *TicTacToeState, msg games.GameMsg) error {
	if len(msg.To) != 1 {
		return errors.New("You must mention exactly one other player")
	}
	state.Fields = []int{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	}
	state.PlayerA = msg.From
	state.PlayerB = msg.To[0]
	state.WhosTurn = 1
	state.Ended = false
	return nil
}

func (t *TicTacToe) renderField(state TicTacToeState) string {
	m := "Field:<br>"
	for i, f := range state.Fields {
		if f == 0 {
			m += intToEmoji(i + 1)
		} else if f == 1 {
			m += "🔵"
		} else if f == 2 {
			m += "🟠"
		}
		if (i+1)%3 == 0 {
			m += "<br>"
		}
	}

	return m
}

// OnMsg implements games.Game.
func (t *TicTacToe) OnMsg(session *models.GameSession, msg games.GameMsg) (interface{}, games.GameReply, error) {
	state := session.Data.(*TicTacToeState)

	// game already ended
	if state.Ended {
		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "The game already ended",
		}, nil
	}

	// init on new game
	if len(state.Fields) != 9 {
		err := t.initState(state, msg)
		if err != nil {
			return state, games.GameReply{
				To:  []string{msg.From},
				Msg: err.Error(),
			}, nil
		}
	}

	// check if it's players turn
	if (state.WhosTurn == 1 && msg.From != state.PlayerA) || (state.WhosTurn == 2 && msg.From != state.PlayerB) {
		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "It's not your turn",
		}, nil
	}

	// apply message to state
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
	// not a valid selection
	if found == 0 {
		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "The message must include a field number between 1 and 9. <br><br>" + t.renderField(*state),
		}, nil
	}
	// field already used
	if state.Fields[found-1] != 0 {
		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "This field is already taken <br><br>" + t.renderField(*state),
		}, nil
	}

	state.Fields[found-1] = state.WhosTurn
	state.WhosTurn = (state.WhosTurn % 2) + 1

	// "print" state to reply
	m := t.renderField(*state)
	m += "<br>"
	actorA, _ := acpub.GetActor(state.PlayerA)
	actorB, _ := acpub.GetActor(state.PlayerB)
	m += "🔵 " + acpub.ActorToLink(actorA) + "<br>"
	m += "🟠 " + acpub.ActorToLink(actorB) + "<br>"

	// check if someone won
	combinations := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
		{1, 5, 9},
		{3, 5, 7},
	}
	winner := 0
	for _, c := range combinations {
		if state.Fields[c[0]-1] == state.Fields[c[1]-1] && state.Fields[c[1]-1] == state.Fields[c[2]-1] {
			if state.Fields[c[0]-1] != 0 {
				winner = state.Fields[c[0]-1]
				break
			}
		}
	}
	if winner != 0 {
		// we have a winner!
		m += "Winner: 🎉🎉🎉 "
		if winner == 1 {
			m += acpub.ActorToLink(actorA)
		} else {
			m += acpub.ActorToLink(actorB)
		}
		m += " 🎉🎉🎉"
		state.Ended = true
	} else {
		// check for draw
		anyFree := false
		for _, f := range state.Fields {
			if f == 0 {
				anyFree = true
				break
			}
		}
		if !anyFree {
			m += "It's a draw."
			state.Ended = true
		} else {
			m += "Its your turn: "

			if state.WhosTurn == 1 {
				m += acpub.ActorToLink(actorA)
			} else {
				m += acpub.ActorToLink(actorB)
			}
		}
	}

	slog.Info("Field message", "msg", m)

	return state, games.GameReply{
		To:  []string{state.PlayerA, state.PlayerB},
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
	return &TicTacToeState{}
}

// Name implements games.Game.
func (TicTacToe) Name() string {
	return "tic-tac-toe"
}
