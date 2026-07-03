package gamelogic

type Player struct {
	Username string
	Team     []Bot
}

type Bot struct {
	Name string
}

type Request struct {
	Command string `json:"command"`
}

func NewPlayer(username string) *Player {
	return &Player{
		Username: username,
	}
}
