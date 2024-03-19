package web

import (
	"net/http"

	"rerere.org/fedi-games/config"
)

func MediaServer() http.Handler {
	cfg := config.GetConfig()
	return http.FileServer(http.Dir(cfg.MediaPath))
}
