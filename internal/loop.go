package internal

import (
	"log/slog"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/domain/models"
	"github.com/H4kor/fedi-games/games"
	"github.com/H4kor/fedi-games/infra"
	"github.com/H4kor/fedi-games/internal/acpub"
	vocab "github.com/go-ap/activitypub"
)

type GameStep struct {
	Sess *models.GameSession
	Game games.Game
	Msg  games.GameMsg
}

type GameEngine struct {
	queue chan GameStep
}

func NewGameEngine() *GameEngine {
	queue := make(chan GameStep, 10)
	engine := &GameEngine{
		queue: queue,
	}
	engine.startProcessor()
	return engine
}

func (engine *GameEngine) startProcessor() {
	go func() {
		for {
			step := <-engine.queue
			engine.process(step.Sess, step.Game, step.Msg)
		}
	}()
}

func (engine *GameEngine) ProcessMsg(sess *models.GameSession, game games.Game, msg games.GameMsg) {
	engine.queue <- GameStep{
		Sess: sess,
		Game: game,
		Msg:  msg,
	}
}

func (engine *GameEngine) process(sess *models.GameSession, game games.Game, msg games.GameMsg) {
	cfg := config.GetConfig()

	// add retrieved message to game session
	sess.MessageIds = append(sess.MessageIds, msg.Id)

	// pass message to game engine
	newState, reply, err := game.OnMsg(sess, msg)

	var note vocab.Note

	if err != nil {
		// on error give the sender a message that there is a problem
		// only sending to sender of message
		slog.Error("Error on Game", "err", err)
		// construct a GameReply for error
		reply := games.GameReply{
			To:  []string{msg.From},
			Msg: "💥 an error occured.",
		}
		// persists the reply
		infra.GetDb().PersistGameReply(sess.GameName, &reply)
		// create Note Object
		note = reply.ToActivityObject(cfg, sess.GameName, msg.Id)
	} else {
		// happy path, Convert To to mentions
		slog.Info("=====================Answer START=========================")
		slog.Info("Answer", "msg", reply.Msg)
		slog.Info("=====================Answer END=========================")
		// persists the reply
		infra.GetDb().PersistGameReply(sess.GameName, &reply)
		// create Note Object
		note = reply.ToActivityObject(cfg, sess.GameName, msg.Id)

	}

	// add note to session messages
	sess.MessageIds = append(sess.MessageIds, note.ID.String())
	// set new state
	sess.Data = newState
	slog.Info("New State", "state", newState)

	// persist game session
	err = infra.GetDb().PersistGameSession(sess)
	if err != nil {
		slog.Error("Error persisting session", "err", err)
	}

	// don't send notes to other services if in localhost mode
	err = acpub.SendNote(sess.GameName, note)
	if err != nil {
		slog.Error("Error sending message", "err", err)
	}
}
