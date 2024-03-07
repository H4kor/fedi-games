package web

import (
	"log/slog"
	"net/http"

	"rerere.org/fedi-games/games"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Serving Index")
	games := make([]games.Game, 0)
	for _, g := range GameMap {
		games = append(games, g)
	}
	err := RenderTemplateWithBase(w, "index", games)
	if err != nil {
		slog.Error("Error rendering index", "err", err)
	}
}
