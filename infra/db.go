package infra

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/H4kor/fedi-games/config"
	"github.com/H4kor/fedi-games/domain/models"
	"github.com/H4kor/fedi-games/games"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqlGameSession struct {
	Id       int64  `db:"id"`
	GameName string `db:"game_name"`
	Data     string `db:"data"`
}

type sqlGameMessage struct {
	Id        string `db:"id"`
	SessionId int64  `db:"session_id"`
}

type sqlGameAttachment struct {
	Url       string
	MediaType string
}

type sqlGameReply struct {
	Id          string `db:"id"`
	GameName    string `db:"game_name"`
	To          string `db:"tos"`
	Msg         string `db:"msg"`
	Attachments string `db:"attachments"`
}

type Database struct {
	db *sqlx.DB
}

var db *Database

func GetDb() *Database {
	if db == nil {
		db = NewDatabase(config.GetConfig().DatabaseUrl)
	}
	return db
}

func NewDatabase(url string) *Database {
	sqlxdb := sqlx.MustOpen("sqlite3", url)

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS game_session (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_name TEXT NOT NULL,
			data TEXT
		);
	`)

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS game_message (
			id TEXT PRIMARY KEY NOT NULL,
			session_id BIGINT NOT NULL,
			FOREIGN KEY(session_id) REFERENCES game_session(id)
		);
	`)

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower TEXT NOT NULL,
			game_name TEXT NOT NULL,
			UNIQUE(follower, game_name) ON CONFLICT IGNORE
		)
	`)

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS game_sent_messages (
			id TEXT PRIMARY KEY,
			game_name TEXT NOT NULL,
			tos TEXT NOT NULL,
			msg TEXT NOT NULL,
			attachments TEXT NOT NULL
		)
	`)

	return &Database{
		db: sqlxdb,
	}
}

func (db *Database) GetGameSessionByMsgId(msgId string, gameName string, gameState interface{}) (*models.GameSession, error) {
	msg := sqlGameMessage{}
	err := db.db.Get(&msg, "SELECT * FROM game_message WHERE id = ?", msgId)
	if err != nil {
		return nil, err
	}

	sess := sqlGameSession{}
	err = db.db.Get(&sess, "SELECT * FROM game_session WHERE id = ? AND game_name = ?", msg.SessionId, gameName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(sess.Data), gameState)
	if err != nil {
		return nil, err
	}

	return &models.GameSession{
		Id:       sess.Id,
		GameName: sess.GameName,
		Data:     gameState,
	}, nil

}

func (db *Database) PersistGameSession(sess *models.GameSession) error {
	data, err := json.Marshal(sess.Data)
	if err != nil {
		return err
	}
	slog.Info("Persisting data", "data", string(data))

	if sess.Id == 0 {
		// new session
		r, err := db.db.Exec("INSERT INTO game_session (game_name, data) VALUES (?, ?)", sess.GameName, string(data))
		if err != nil {
			return err
		}

		id, err := r.LastInsertId()
		if err != nil {
			return err
		}
		sess.Id = id
	} else {
		// update session
		_, err = db.db.Exec("UPDATE game_session SET data = ? WHERE id = ?", string(data), sess.Id)
		if err != nil {
			return err
		}
	}

	// persist new messages
	for _, m := range sess.MessageIds {
		_, err = db.db.Exec("INSERT INTO game_message (id, session_id) VALUES (?, ?) ON CONFLICT DO NOTHING", m, sess.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) AddFollower(gameName string, follower string) error {
	_, err := db.db.Exec("INSERT INTO followers (game_name, follower) VALUES (?, ?)", gameName, follower)
	return err
}

func (db *Database) RemoveFollower(gameName string, follower string) error {
	_, err := db.db.Exec("DELETE FROM followers WHERE game_name = ? AND follower = ?", gameName, follower)
	return err
}

func (db *Database) ListFollowers(gameName string) ([]string, error) {
	slog.Info("Getting followers for", "gameName", gameName)
	followers := []string{}
	err := db.db.Select(&followers, "SELECT follower FROM followers WHERE game_name = ?", gameName)
	if err != nil {
		return []string{}, err
	}

	return followers, err
}

func (db *Database) PersistGameReply(gameName string, reply *games.GameReply) error {
	if reply.Id == "" {
		id, _ := uuid.NewRandom()
		reply.Id = id.String()
		tos, _ := json.Marshal(reply.To)
		sqlAtt := make([]sqlGameAttachment, 0, len(reply.Attachments))
		for _, a := range reply.Attachments {
			sqlAtt = append(sqlAtt, sqlGameAttachment{
				Url:       a.Url,
				MediaType: a.MediaType,
			})
		}
		attachments, _ := json.Marshal(sqlAtt)

		_, err := db.db.Exec(`INSERT INTO game_sent_messages (id, game_name, tos, msg, attachments) VALUES (?, ?, ?, ?, ?)`,
			id, gameName, tos, reply.Msg, attachments,
		)
		return err
	} else {
		return errors.New("reply already has an id, will not be persisted")
	}
}

func (db *Database) RetrieveGameReply(gameName string, id string) (*games.GameReply, error) {
	slog.Info("Getting followers for", "gameName", gameName)
	sqlReply := sqlGameReply{}
	err := db.db.Get(&sqlReply, "SELECT * FROM game_sent_messages WHERE game_name = ? AND id = ?", gameName, id)
	if err != nil {
		return &games.GameReply{}, err
	}
	attachments := make([]games.GameAttachment, 0, len(sqlReply.Attachments))
	aa := []sqlGameAttachment{}
	json.Unmarshal([]byte(sqlReply.Attachments), &aa)
	for _, a := range aa {
		attachments = append(attachments, games.GameAttachment{
			Url:       a.Url,
			MediaType: a.MediaType,
		})
	}
	to := []string{}
	json.Unmarshal([]byte(sqlReply.To), &to)

	reply := games.GameReply{
		Id:          sqlReply.Id,
		To:          to,
		Msg:         sqlReply.Msg,
		Attachments: attachments,
	}

	return &reply, err
}
