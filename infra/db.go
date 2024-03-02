package infra

import (
	"encoding/json"

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
			id BIGINT PRIMARY KEY,
			game_name TEXT,
			state TEXT
		);
	`)

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS game_message (
			id TEXT PRIMARY KEY,
			session_id BIGINT,
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
