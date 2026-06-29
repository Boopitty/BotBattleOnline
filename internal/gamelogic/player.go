package gamelogic

type Player struct {
	Username string
	Team     []bot
}

func NewPlayer(username string) *Player {
	return &Player{
		Username: username,
	}
}
