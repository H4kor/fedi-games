package internal

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	vocab "github.com/go-ap/activitypub"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/domain/models"
	"rerere.org/fedi-games/games"
	"rerere.org/fedi-games/infra"
	"rerere.org/fedi-games/internal/acpub"
)

func ProcessMsg(sess *models.GameSession, game games.Game, msg games.GameMsg) {
	cfg := config.GetConfig()

	// add retrieved message to game session
	sess.MessageIds = append(sess.MessageIds, msg.Id)

	// pass message to game engine
	newState, ret, err := game.OnMsg(sess, msg)

	// contruct params
	to := vocab.ItemCollection{}
	mentions := vocab.ItemCollection{}
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
		for _, t := range ret.To {
			to = append(to, vocab.ID(t))
			mention := vocab.MentionNew(
				vocab.ID(t),
			)
			mention.Href = vocab.ID(t)
			mentions = append(mentions, mention)
			msgStr = ret.Msg
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
