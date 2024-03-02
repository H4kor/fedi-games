package main

import (
	"net/http"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/web"
)

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

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /.well-known/webfinger", web.WebfingerHandler)
	mux.HandleFunc("GET /games/{game}", gameHandler)
	mux.HandleFunc("POST /games/{game}/inbox", web.InboxHandler)

	println("Starting server on port 4040")
	err := http.ListenAndServe(":4040", mux)
	if err != nil {
		panic(err)
	}

}
