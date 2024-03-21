package web

import (
	"net/http"

	"github.com/H4kor/fedi-games/config"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
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
