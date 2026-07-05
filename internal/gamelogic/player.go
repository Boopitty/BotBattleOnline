package gamelogic

func NewPlayer() *Player {
	return &Player{
		Username: "",
		Team:     make(map[int]Bot),
	}
}
