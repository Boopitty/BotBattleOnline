package gamelogic

type Player struct {
	Username string
	Team     []Bot
}

func NewPlayer(username string) *Player {
	return &Player{
		Username: username,
	}
}
