package games

type Game interface {
	OnMsg(GameMsg) (GameReply, error)
}

type GameMsg struct {
	Id      string
	From    string
	To      []string
	Msg     string
	ReplyTo *string
}

type GameReply struct {
	To  []string
	Msg string
}
