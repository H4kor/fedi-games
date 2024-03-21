package web

import (
	"embed"
	"net/http"
)

//go:embed static
var staticFs embed.FS

func (server *FediGamesServer) StaticServer() http.Handler {
	return http.FileServer(http.FS(staticFs))
}
