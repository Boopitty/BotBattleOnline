package gamelogic

import (
	"sync"
)

type GameState struct {
	Players map[int]Player `json:"players"`
	Paused  bool           `json:"paused"`
	Mu      *sync.RWMutex
}

func NewGameState() *GameState {
	return &GameState{
		Players: make(map[int]Player),
		Paused:  false,
		Mu:      &sync.RWMutex{},
	}
}
