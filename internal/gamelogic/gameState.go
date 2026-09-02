package gamelogic

import (
	"sync"

	"github.com/gorilla/websocket"
)

type GameState struct {
	Players    map[*websocket.Conn]Player `json:"players"`
	Spectators map[*websocket.Conn]Player `json:"spectators"`
	ActiveBots []Bot                      `json:"active_bots"`
	Turn       int                        `json:"turn"`
	Mu         *sync.RWMutex
}

func NewGameState() *GameState {
	return &GameState{
		Players:    make(map[*websocket.Conn]Player, 2),
		Spectators: make(map[*websocket.Conn]Player),
		ActiveBots: make([]Bot, 0, 4),
		Turn:       0,
		Mu:         &sync.RWMutex{},
	}
}
