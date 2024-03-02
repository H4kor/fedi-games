package main

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/games"
	tictactoe "rerere.org/fedi-games/games/tic-tac-toe"
	"rerere.org/fedi-games/internal/acpub"
	"rerere.org/fedi-games/internal/html"
)

type WebfingerResponse struct {
	Subject string          `json:"subject"`
	Aliases []string        `json:"aliases"`
	Links   []WebfingerLink `json:"links"`
}

type WebfingerLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

var gamesMap = map[string]games.Game{
	"tic-tac-toe": tictactoe.NewTicTacToeGame(),
}

func writeJson(w http.ResponseWriter, data interface{}) error {
	s, _ := json.Marshal(data)
	w.Write(s)
	return nil
}

func webfingerHandler(w http.ResponseWriter, r *http.Request) {
	host := config.GetConfig().Host

	resource := r.URL.Query().Get("resource")
	if !strings.HasPrefix(resource, "acct:") {
		println("error: must start with acct:")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	parts := strings.Split(resource[5:], "@")
	if len(parts) != 2 {
		println("error: must have @")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	req_name := parts[0]
	req_host := parts[1]

	if req_host != host {
		println("error: not host")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, ok := gamesMap[req_name]
	if !ok {
		println("error: unknown game")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cfg := config.GetConfig()

	webfinger := WebfingerResponse{
		Subject: resource,

		Links: []WebfingerLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: cfg.FullUrl() + "/games/" + req_name,
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	writeJson(w, webfinger)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	game := r.PathValue("game")

	cfg := config.GetConfig()

	actor := vocab.ServiceNew(vocab.IRI(cfg.FullUrl() + "/games/" + game))
	actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(game)}}
	actor.Inbox = vocab.IRI(cfg.FullUrl() + "/games/" + game + "/inbox")
	// actor.Outbox = vocab.IRI(config.FullUrl() + "/games/" + game + "/outbox")
	// actor.Followers = vocab.IRI(config.FullUrl() + "/games/" + game + "/followers")
	actor.PublicKey = vocab.PublicKey{
		ID:           vocab.ID(cfg.FullUrl() + "/games/" + game + "#main-key"),
		Owner:        vocab.IRI(cfg.FullUrl() + "/games/" + game),
		PublicKeyPem: cfg.PublicKeyPem,
	}
	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
		jsonld.IRI(vocab.SecurityContextURI),
	).Marshal(actor)
	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}

func inboxHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")

	game, ok := gamesMap[gameName]
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
				participants = append(participants, r.GetLink().String())
			}

			var replyTo *string
			if o.InReplyTo != nil {
				r := o.InReplyTo.GetLink().String()
				replyTo = &r
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

			go handleGameStep(gameName, game, gameMsg)

			slog.Info("Done")

			return nil
		})
	})
	if err != nil {
		slog.Error("Error on activity", "err", err)
	}
}

func handleGameStep(gameName string, game games.Game, msg games.GameMsg) {
	ret, err := game.OnMsg(msg)
	if err != nil {
		slog.Error("Error on Game", "err", err)
	}

	cfg := config.GetConfig()

	to := vocab.ItemCollection{}
	mentions := vocab.ItemCollection{}
	for _, t := range ret.To {
		to = append(to, vocab.ID(t))
		mentions = append(mentions, vocab.MentionNew(
			vocab.ID(t),
		))
	}

	note := vocab.Note{
		ID:           vocab.ID(cfg.FullUrl() + "/games/" + gameName + "/" + strconv.FormatInt(time.Now().Unix(), 10)),
		Type:         "Note",
		InReplyTo:    vocab.ID(msg.Id),
		To:           to,
		Published:    time.Now().UTC(),
		AttributedTo: vocab.ID(cfg.FullUrl() + "/games/" + gameName),
		Tag:          mentions,
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(ret.Msg)},
		},
	}

	if err != nil {
		slog.Error("Error Marshalling", "err", err)
	}

	err = acpub.SendNote(gameName, note)
	if err != nil {
		slog.Error("Error sending message", "err", err)
	}

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /.well-known/webfinger", webfingerHandler)
	mux.HandleFunc("GET /games/{game}", gameHandler)
	mux.HandleFunc("POST /games/{game}/inbox", inboxHandler)

	println("Starting server on port 4040")
	err := http.ListenAndServe(":4040", mux)
	if err != nil {
		panic(err)
	}

}
