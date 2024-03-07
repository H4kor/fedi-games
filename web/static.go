package web

import (
	"embed"
	"net/http"
)

//go:embed static
var staticFs embed.FS

func StaticServer() http.Handler {
	return http.FileServer(http.FS(staticFs))
}
