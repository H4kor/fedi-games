package web

import (
	"net/http"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"rerere.org/fedi-games/config"
)

func (server *FediGamesServer) FollowingHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	game, ok := server.games[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cfg := config.GetConfig()

	following := vocab.Collection{}

	for _, g := range server.games {
		following.Append(gameUrl(g))
	}
	following.TotalItems = uint(len(server.games))
	following.ID = vocab.IRI(cfg.FullUrl() + "/games/" + game.Name() + "/following")
	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(following)

	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}
