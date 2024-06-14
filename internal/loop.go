package internal

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/domain/models"
	"github.com/H4kor/fedi-games/games"
	"github.com/H4kor/fedi-games/infra"
	"github.com/H4kor/fedi-games/internal/acpub"
	vocab "github.com/go-ap/activitypub"
)

type GameStep struct {
	Sess *models.GameSession
	Game games.Game
	Msg  games.GameMsg
}

type GameEngine struct {
	queue chan GameStep
}

func NewGameEngine() *GameEngine {
	queue := make(chan GameStep, 10)
	engine := &GameEngine{
		queue: queue,
	}
	engine.startProcessor()
	return engine
}

func (engine *GameEngine) startProcessor() {
	go func() {
		for {
			step := <-engine.queue
			engine.process(step.Sess, step.Game, step.Msg)
		}
	}()
}

func (engine *GameEngine) ProcessMsg(sess *models.GameSession, game games.Game, msg games.GameMsg) {
	engine.queue <- GameStep{
		Sess: sess,
		Game: game,
		Msg:  msg,
	}
}

func (engine *GameEngine) process(sess *models.GameSession, game games.Game, msg games.GameMsg) {
	cfg := config.GetConfig()

	// add retrieved message to game session
	sess.MessageIds = append(sess.MessageIds, msg.Id)

	// pass message to game engine
	newState, reply, err := game.OnMsg(sess, msg)

	// contruct params
	to := vocab.ItemCollection{}
	mentions := vocab.ItemCollection{}
	attachments := vocab.ItemCollection{}
	msgStr := ""

	if err != nil {
		// on error give the sender a message that there is a problem
		// only sending to sender of message
		slog.Error("Error on Game", "err", err)
		to = append(to, vocab.ID(msg.From))
		mention := vocab.MentionNew(
			vocab.ID(msg.From),
		)
		mention.Href = vocab.ID(msg.From)
		mentions = append(mentions, mention)
		msgStr = "💥 an error occured."
	} else {
		// happy path, Convert To to mentions
		slog.Info("=====================Answer START=========================")
		slog.Info("Answer", "msg", reply.Msg)
		slog.Info("=====================Answer END=========================")
		for _, t := range reply.To {
			to = append(to, vocab.ID(t))
			mention := vocab.MentionNew(
				vocab.ID(t),
			)
			mention.Href = vocab.ID(t)
			mentions = append(mentions, mention)
			msgStr = reply.Msg
		}
		for _, a := range reply.Attachments {
			att := vocab.Document{
				Type:      vocab.DocumentType,
				MediaType: vocab.MimeType(a.MediaType),
				URL:       vocab.ID(a.Url),
			}
			attachments = append(attachments, att)
		}
	}

	// Construct Activitystream Note
	note := vocab.Note{
		ID:           vocab.ID(cfg.FullUrl() + "/games/" + sess.GameName + "/" + strconv.FormatInt(time.Now().Unix(), 10)),
		Type:         "Note",
		InReplyTo:    vocab.ID(msg.Id),
		To:           to,
		Published:    time.Now().UTC(),
		AttributedTo: vocab.ID(cfg.FullUrl() + "/games/" + sess.GameName),
		Tag:          mentions,
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(msgStr)},
		},
		Attachment: attachments,
	}
	// add note to session messages
	sess.MessageIds = append(sess.MessageIds, note.ID.String())
	// set new state
	sess.Data = newState
	slog.Info("New State", "state", newState)

	// persist game session
	err = infra.GetDb().PersistGameSession(sess)
	if err != nil {
		slog.Error("Error persisting session", "err", err)
	}

	// don't send notes to other services if in localhost mode
	if !strings.Contains(cfg.FullUrl(), "localhost") {
		err = acpub.SendNote(sess.GameName, note)
		if err != nil {
			slog.Error("Error sending message", "err", err)
		}
	}

}
