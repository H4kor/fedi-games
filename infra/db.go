package infra

import (
	"encoding/json"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"rerere.org/fedi-games/config"
	"rerere.org/fedi-games/domain/models"
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

type Database struct {
	db *sqlx.DB
}

var db *Database

func GetDb() *Database {
	if db == nil {
		db = NewDatabase()
	}
	return db
}

func NewDatabase() *Database {
	sqlxdb := sqlx.MustOpen("sqlite3", config.GetConfig().DatabaseUrl)

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
