package web

import (
	"net/http"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

func OutboxHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")

	_, ok := GameMap[gameName]
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
