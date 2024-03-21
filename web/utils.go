package web

import (
	"encoding/json"
	"net/http"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/games"
	vocab "github.com/go-ap/activitypub"
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
