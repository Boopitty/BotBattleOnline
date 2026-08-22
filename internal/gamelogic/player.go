package gamelogic

type Player struct {
	Username string      `json:"username"`
	Team     map[int]Bot `json:"team"`
}

func NewPlayer(username string) *Player {
	return &Player{
		Username: username,
		Team:     make(map[int]Bot),
	}
}
