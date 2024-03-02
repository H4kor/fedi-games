package models

type GameSession struct {
	Id         int64
	GameName   string
	Data       interface{}
	MessageIds []string
}
