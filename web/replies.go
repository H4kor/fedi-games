package web

import (
	"net/http"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/infra"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

func (server *FediGamesServer) ReplyHandler(w http.ResponseWriter, r *http.Request) {
	gameName := r.PathValue("game")
	replyId := r.PathValue("game")
	_, ok := server.games[gameName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cfg := config.GetConfig()

	reply, err := infra.GetDb().RetrieveGameReply(gameName, replyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(reply.ToActivityObject(cfg, gameName, ""))

	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(data)
}
