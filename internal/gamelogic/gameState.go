package gamelogic

import (
	"sync"
)

func NewGameState() *GameState {
	return &GameState{
		Players: make(map[int]Player),
		Paused:  false,
		Mu:      &sync.RWMutex{},
	}
}
