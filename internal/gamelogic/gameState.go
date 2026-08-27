package gamelogic

import (
	"sync"
)

type GameState struct {
	Players    map[int]Player `json:"players"`
	Spectators map[int]Player `json:"spectators"`
	ActiveBots []Bot          `json:"active_bots"`
	Turn       int            `json:"turn"`
	Mu         *sync.RWMutex
}

func NewGameState() *GameState {
	return &GameState{
		Players:    make(map[int]Player),
		Spectators: make(map[int]Player),
		ActiveBots: make([]Bot, 4),
		Turn:       0,
		Mu:         &sync.RWMutex{},
	}
}
