package games

import (
	"time"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/domain/models"
	vocab "github.com/go-ap/activitypub"
)

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
	Id          string
	To          []string
	Msg         string
	Attachments []GameAttachment
}

func (r *GameReply) ToActivityObject(cfg config.Config, gameName string, inReplyTo string) vocab.Note {
	to := vocab.ItemCollection{}
	mentions := vocab.ItemCollection{}
	attachments := vocab.ItemCollection{}
	msgStr := r.Msg

	for _, t := range r.To {
		to = append(to, vocab.ID(t))
		mention := vocab.MentionNew(
			vocab.ID(t),
		)
		mention.Href = vocab.ID(t)
		mentions = append(mentions, mention)
	}
	for _, a := range r.Attachments {
		att := vocab.Document{
			Type:      vocab.DocumentType,
			MediaType: vocab.MimeType(a.MediaType),
			URL:       vocab.ID(a.Url),
		}
		attachments = append(attachments, att)
	}

	// Construct Activitystream Note
	return vocab.Note{
		ID:           vocab.ID(cfg.FullUrl() + "/games/" + gameName + "/replies/" + r.Id),
		Type:         "Note",
		InReplyTo:    vocab.ID(inReplyTo),
		To:           to,
		Published:    time.Now().UTC(),
		AttributedTo: vocab.ID(cfg.FullUrl() + "/games/" + gameName),
		Tag:          mentions,
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(msgStr)},
		},
		Attachment: attachments,
	}
}
