package web

import (
	"net/http"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

func (server *FediGamesServer) FollowersHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")

	_, ok := server.games[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	following := vocab.Collection{}

	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(following)

	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}
