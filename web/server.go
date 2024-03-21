package web

import (
	"net/http"

	"rerere.org/fedi-games/games"
	"rerere.org/fedi-games/internal"
)

type FediGamesServer struct {
	engine *internal.GameEngine
	games  map[string]games.Game
}

func NewFediGamesServer(engine *internal.GameEngine, gamesList []games.Game) FediGamesServer {
	gamesMap := map[string]games.Game{}
	for _, g := range gamesList {
		gamesMap[g.Name()] = g
	}
	return FediGamesServer{
		engine: engine,
		games:  gamesMap,
	}
}

func (server *FediGamesServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/webfinger", server.WebfingerHandler)
	mux.HandleFunc("GET /games/{game}", server.GameHandler)
	mux.HandleFunc("POST /games/{game}/inbox", server.InboxHandler)
	mux.HandleFunc("GET /games/{game}/outbox", server.OutboxHandler)
	mux.HandleFunc("GET /games/{game}/following", server.FollowingHandler)
	mux.HandleFunc("GET /games/{game}/followers", server.FollowersHandler)
	mux.HandleFunc("GET /games/{game}/avatar.png", server.AvatarHandler)
	mux.Handle("GET /static/", server.StaticServer())
	mux.Handle("GET /media/", http.StripPrefix("/media", server.MediaServer()))
	mux.HandleFunc("GET /{$}", server.IndexHandler)

	println("Starting server on port 4040")
	return http.ListenAndServe(":4040", mux)
}
