package web

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/H4kor/fedi-games/config"
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

func (server *FediGamesServer) WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	host := config.GetConfig().Host

	resource := r.URL.Query().Get("resource")
	if !strings.HasPrefix(resource, "acct:") {
		slog.Error("resource must start with acct:", "resource", resource)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	parts := strings.Split(resource[5:], "@")
	if len(parts) != 2 {
		slog.Error("resource must have @", "resource", resource)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	req_name := parts[0]
	req_host := parts[1]

	if req_host != host {
		slog.Error("not host", "req_host", req_host, "host", host)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, ok := server.games[req_name]
	if !ok {
		slog.Error("unknown game", "req_name", req_name)
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
