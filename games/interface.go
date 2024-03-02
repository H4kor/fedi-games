package games

import "rerere.org/fedi-games/domain/models"

type Game interface {
	OnMsg(*models.GameSession, GameMsg) (GameReply, error)
	NewState() interface{}
	Name() string
}

type GameMsg struct {
	Id      string
	From    string
	To      []string
	Msg     string
	ReplyTo *string
}

type GameReply struct {
	To  []string
	Msg string
}
