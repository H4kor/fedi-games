package web

import (
	"log/slog"
	"net/http"

	"rerere.org/fedi-games/games"
)

func (server *FediGamesServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Serving Index")
	games := make([]games.Game, 0)
	for _, g := range server.games {
		games = append(games, g)
	}
	err := RenderTemplateWithBase(w, "index", games)
	if err != nil {
		slog.Error("Error rendering index", "err", err)
	}
}
