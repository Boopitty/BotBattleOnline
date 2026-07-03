package gamelogic

type Player struct {
	Username string      `json:"username"`
	Team     map[int]Bot `json:"team"`
}

func NewPlayer() *Player {
	return &Player{
		Username: "",
		Team:     make(map[int]Bot),
	}
}
