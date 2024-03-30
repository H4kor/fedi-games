package web

import (
	"log/slog"
	"net/http"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/infra"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

func (server *FediGamesServer) FollowersHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	game, ok := server.games[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cfg := config.GetConfig()

	followers := vocab.Collection{}

	fs, err := infra.GetDb().ListFollowers(gameName)
	if err != nil {
		slog.Error("error on followers", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Info("followers", "fs", fs)
	for _, f := range fs {
		followers.Append(vocab.IRI(f))
	}
	followers.TotalItems = uint(len(fs))
	followers.ID = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/followers")
	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(followers)

	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}
