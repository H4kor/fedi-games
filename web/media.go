package web

import (
	"net/http"

	"github.com/H4kor/fedi-games/config"
)

func (server *FediGamesServer) MediaServer() http.Handler {
	cfg := config.GetConfig()
	return http.FileServer(http.Dir(cfg.MediaPath))
}
