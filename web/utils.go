package web

import (
	"encoding/json"
	"net/http"

	vocab "github.com/go-ap/activitypub"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/games"
)

func writeJson(w http.ResponseWriter, data interface{}) error {
	s, _ := json.Marshal(data)
	w.Write(s)
	return nil
}

func gameUrl(game games.Game) vocab.IRI {
	cfg := config.GetConfig()
	return vocab.IRI(cfg.FullUrl() + "/games/" + game.Name())
}
