package web

import (
	"net/http"

	"rerere.org/fedi-games/config"
)

func (server *FediGamesServer) MediaServer() http.Handler {
	cfg := config.GetConfig()
	return http.FileServer(http.Dir(cfg.MediaPath))
}
