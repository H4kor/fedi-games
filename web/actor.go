package web

import (
	"log/slog"
	"net/http"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/games"
)

type htmlData struct {
	Game games.Game
	Cfg  config.Config
}

func GameHandlerHtml(w http.ResponseWriter, r *http.Request, game games.Game) {
	cfg := config.GetConfig()
	err := RenderTemplateWithBase(w, "game", htmlData{
		Game: game,
		Cfg:  cfg,
	})
	if err != nil {
		slog.Error("Error rendering index", "err", err)
	}
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	game, ok := GameMap[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		GameHandlerHtml(w, r, game)
		return
	}

	cfg := config.GetConfig()

	actor := vocab.ServiceNew(vocab.IRI(cfg.FullUrl() + "/games/" + game.Name()))
	actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(game.Name())}}
	actor.Inbox = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/inbox")
	actor.Following = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/following")
	actor.Outbox = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/outbox")
	// actor.Followers = vocab.IRI(config.FullUrl() + "/games/" + game + "/followers")
	actor.PublicKey = vocab.PublicKey{
		ID:           vocab.ID(cfg.FullUrl() + "/games/" + game.Name() + "#main-key"),
		Owner:        vocab.IRI(cfg.FullUrl() + "/games/" + game.Name()),
		PublicKeyPem: cfg.PublicKeyPem,
	}

	actor.Name = vocab.NaturalLanguageValues{{Value: vocab.Content(game.Name())}}
	actor.Icon = vocab.Image{
		Type:      vocab.ImageType,
		MediaType: "image/png",
		URL:       vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/avatar.png"),
	}
	actor.Summary = vocab.NaturalLanguageValues{{Value: vocab.Content(game.Summary())}}

	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
		jsonld.IRI(vocab.SecurityContextURI),
	).Marshal(actor)
	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}
