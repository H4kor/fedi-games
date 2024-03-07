package web

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/domain/models"
	"rerere.org/fedi-games/games"
	"rerere.org/fedi-games/infra"
	"rerere.org/fedi-games/internal"
	"rerere.org/fedi-games/internal/html"
)

type InboxHandler struct {
	engine *internal.GameEngine
}

func NewInboxHandler(engine *internal.GameEngine) http.Handler {
	return &InboxHandler{
		engine: engine,
	}
}

// ServeHTTP implements http.Handler.
func (i *InboxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	cfg := config.GetConfig()

	game, ok := GameMap[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Couldn't read body in inbox")
	}
	slog.Info("Retrieved in inbox", "body", string(body))
	data, err := vocab.UnmarshalJSON(body)
	if err != nil {
		slog.Error("Couldn't unmarshal body", "err", err)
	}
	slog.Info("Data retrieved", "data", data)

	err = vocab.OnActivity(data, func(act *vocab.Activity) error {
		slog.Info("activity retrieved", "activity", act)
		if act.Type != "Create" {
			return errors.New("only create activities are supported")
		}

		return vocab.OnObject(act.Object, func(o *vocab.Object) error {
			slog.Info("object retrieved", "object", o)
			plain := html.GetTextFromHtml(o.Content.String())
			recipients := o.Recipients()
			sender := o.AttributedTo.GetLink().String()
			participants := []string{}
			for _, r := range recipients {
				// filter out all actors on this server fromt the participants list
				if strings.Contains(r.GetLink().String(), cfg.Host) {
					continue
				}
				// filter out special @s ( e.g. https://www.w3.org/ns/activitystreams#Public )
				if strings.Contains(r.GetLink().String(), "https://www.w3.org/ns/activitystreams") {
					continue
				}
				// filter out sender
				if strings.Contains(r.GetLink().String(), sender) {
					continue
				}

				participants = append(participants, r.GetLink().String())
			}

			state := game.NewState()
			var sess *models.GameSession

			var replyTo *string
			if o.InReplyTo != nil {
				r := o.InReplyTo.GetLink().String()
				replyTo = &r
				sess, err = infra.GetDb().GetGameSessionByMsgId(r, gameName, state)
				if err != nil {
					slog.Error("Error getting game state", "err", err)
					sess = &models.GameSession{
						GameName: gameName,
						Data:     state,
					}
				}
			} else {
				sess = &models.GameSession{
					GameName: gameName,
					Data:     state,
				}
			}

			gameMsg := games.GameMsg{
				Id:      o.ID.String(),
				From:    sender,
				To:      participants,
				Msg:     plain,
				ReplyTo: replyTo,
			}

			slog.Info("Content of object", "content", o.Content.String(), "plain", plain)
			slog.Info("Game Message", "msg", gameMsg)

			i.engine.ProcessMsg(sess, game, gameMsg)

			slog.Info("Done")

			return nil
		})
	})
	if err != nil {
		slog.Error("Error on activity", "err", err)
	}
}
