package gamelogic

import (
	"sync"
)

type GameState struct {
	Players struct {
		Player1 Player `json:"player1"`
		Player2 Player `json:"player2"`
	} `json:"players"`
	Spectators []Player `json:"spectators"`
	ActiveBots []Bot    `json:"active_bots"`
	Turn       int      `json:"turn"`
	Mu         *sync.RWMutex
}

func NewGameState() *GameState {
	return &GameState{
		Players: struct {
			Player1 Player `json:"player1"`
			Player2 Player `json:"player2"`
		}{
			Player1: Player{},
			Player2: Player{},
		},
		Spectators: make([]Player, 0),
		ActiveBots: make([]Bot, 0, 4),
		Turn:       0,
		Mu:         &sync.RWMutex{},
	}
}
