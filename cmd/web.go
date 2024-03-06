package main

import (
	"net/http"

	"rerere.org/fedi-games/web"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /.well-known/webfinger", web.WebfingerHandler)
	mux.HandleFunc("GET /games/{game}", web.GameHandler)
	mux.HandleFunc("POST /games/{game}/inbox", web.InboxHandler)
	mux.HandleFunc("GET /games/{game}/outbox", web.OutboxHandler)
	mux.HandleFunc("GET /games/{game}/following", web.FollowingHandler)
	mux.HandleFunc("GET /games/{game}/followers", web.FollowersHandler)
	mux.HandleFunc("GET /games/{game}/avatar.png", web.AvatarHandler)

	println("Starting server on port 4040")
	err := http.ListenAndServe(":4040", mux)
	if err != nil {
		panic(err)
	}

}
