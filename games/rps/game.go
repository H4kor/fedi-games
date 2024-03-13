package rps

import (
	"math/rand"
	"strings"

	"rerere.org/fedi-games/domain/models"
	"rerere.org/fedi-games/games"
)

type RockPaperScissor struct{}

type RockPaperScissorState struct{}

// Name implements games.Game.
func (r RockPaperScissor) Name() string {
	return "rps"
}

// NewState implements games.Game.
func (r *RockPaperScissor) NewState() interface{} {
	return &RockPaperScissorState{}
}

// OnMsg implements games.Game.
func (r *RockPaperScissor) OnMsg(sess *models.GameSession, msg games.GameMsg) (interface{}, games.GameReply, error) {
	parts := strings.Split(msg.Msg, " ")
	selection := ""
	for _, p := range parts {
		p = strings.ToLower(p)
		if p == "rock" || p == "🪨" {
			selection = "rock"
		}
		if p == "paper" || p == "📄" {
			selection = "paper"
		}
		if p == "scissor" || p == "scissors" || p == "✂️" {
			selection = "scissor"
		}
		if selection != "" {
			break
		}
	}

	if selection == "" {
		return sess, games.GameReply{
			To:  []string{msg.From},
			Msg: "Your message must contain 'rock', 'paper' or 'scissors'.",
		}, nil
	}

	var bot string
	switch rand.Intn(3) {
	case 0:
		bot = "rock"
		break
	case 1:
		bot = "paper"
		break
	default:
		bot = "scissor"
	}

	result := "draw"
	if bot == "rock" {
		if selection == "paper" {
			result = "player"
		}
		if selection == "scissor" {
			result = "bot"
		}
	}
	if bot == "paper" {
		if selection == "scissor" {
			result = "player"
		}
		if selection == "rock" {
			result = "bot"
		}
	}
	if bot == "scissor" {
		if selection == "rock" {
			result = "player"
		}
		if selection == "paper" {
			result = "bot"
		}
	}

	m := "I choose: "
	if bot == "rock" {
		m += "🪨"
	}
	if bot == "paper" {
		m += "📄"
	}
	if bot == "scissor" {
		m += "✂️"
	}
	m += "<br><br>"
	if result == "draw" {
		m += "DRAW!"
	}
	if result == "player" {
		m += "YOU WIN! 🎉"
	}
	if result == "bot" {
		m += "I WIN! 😈"
	}

	return sess, games.GameReply{
		To:  []string{msg.From},
		Msg: m,
	}, nil
}

// Summary implements games.Game.
func (r *RockPaperScissor) Summary() string {
	return "Just 'Pock Paper Scissors!'. Write me your choice to start a game.<br>I promise I won't cheat."
}

// Example implements games.Game.
func (r *RockPaperScissor) Example() string {
	return "rock"
}

func NewRockPaperScissorGame() games.Game {
	return &RockPaperScissor{}
}
