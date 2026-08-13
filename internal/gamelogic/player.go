package gamelogic

func NewPlayer(username string) *Player {
	return &Player{
		Username: username,
		Team:     make(map[int]Bot),
	}
}
