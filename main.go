package main

import (
	"encoding/json"
	"net/http"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
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

type Config struct {
	Host         string
	Protocol     string
	PublicKeyPem string
}

func (c *Config) FullUrl() string {
	return c.Protocol + "://" + c.Host
}

func config() Config {
	return Config{
		Host:         "localhost:4040",
		Protocol:     "http",
		PublicKeyPem: "TBD",
	}
}

var games = map[string]interface{}{
	"tic-tac-toe": 1,
}

func writeJson(w http.ResponseWriter, data interface{}) error {
	s, _ := json.Marshal(data)
	w.Write(s)
	return nil
}

func webfingerHandler(w http.ResponseWriter, r *http.Request) {
	host := config().Host

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

	_, ok := games[req_name]
	if !ok {
		println("error: unknown game")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	webfinger := WebfingerResponse{
		Subject: resource,

		Links: []WebfingerLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: config().Protocol + "://" + config().Host + "/games/" + req_name,
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	writeJson(w, webfinger)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	game := r.PathValue("game")

	config := config()

	actor := vocab.PersonNew(vocab.IRI(config.FullUrl() + "/games/" + game))
	actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(game)}}
	actor.Inbox = vocab.IRI(config.FullUrl() + "/games/" + game + "/inbox")
	actor.Outbox = vocab.IRI(config.FullUrl() + "/games/" + game + "/outbox")
	actor.Followers = vocab.IRI(config.FullUrl() + "/games/" + game + "/followers")
	actor.PublicKey = vocab.PublicKey{
		ID:           vocab.ID(config.FullUrl() + "/games/" + game + "/actor#main-key"),
		Owner:        vocab.IRI(config.FullUrl() + "/games/" + game),
		PublicKeyPem: config.PublicKeyPem,
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

	mux.HandleFunc("GET /.well-known/webfinger", webfingerHandler)
	mux.HandleFunc("GET /games/{game}", gameHandler)

	println("Starting server on port 4040")
	err := http.ListenAndServe(":4040", mux)
	if err != nil {
		panic(err)
	}

}
