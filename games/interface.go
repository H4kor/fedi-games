package games

import "rerere.org/fedi-games/domain/models"

type Game interface {
	OnMsg(*models.GameSession, GameMsg) (interface{}, GameReply, error)
	NewState() interface{}
	Name() string
	Summary() string
	Example() string
}

type GameMsg struct {
	Id      string
	From    string
	To      []string
	Msg     string
	ReplyTo *string
}

type GameAttachment struct {
	Url       string
	MediaType string
}

type GameReply struct {
	To          []string
	Msg         string
	Attachments []GameAttachment
}
