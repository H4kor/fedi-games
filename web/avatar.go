package web

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

//go:embed avatars
var avatars embed.FS

func AvatarHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	_, ok := GameMap[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}

	bytes, err := fs.ReadFile(avatars, "avatars/"+gameName+".png")
	if err != nil {
		slog.Error("Error loading avatar", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "image/png")
	w.Write(bytes)

}
