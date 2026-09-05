package gamelogic

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn `json:"-"`
	Username string          `json:"username"`
	Team     map[int]Bot     `json:"team"`
}

func NewPlayer(username string) *Player {
	return &Player{
		Username: username,
		Team:     make(map[int]Bot),
	}
}
