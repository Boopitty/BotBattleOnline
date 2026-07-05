package gamelogic

import (
	"sync"
)

// Request represents a command request from a player.
type Request struct {
	User    Player `json:"user"`
	Command string `json:"command"`
	Message string `json:"message"`
}

// Used when creating a new gamestate
type Player struct {
	Username string      `json:"username"`
	Team     map[int]Bot `json:"team"`
}

// Store information about bots
type Bot struct {
	Name    string `json:"name"`
	Health  int    `json:"health"`
	Power   int    `json:"power"`
	Defense int    `json:"defense"`
	Huge    bool   `json:"huge"`
	IsAlive bool   `json:"is_alive"`
}

// Used to share and store the gamestate of a multiplayer game
type GameState struct {
	Players map[int]Player `json:"players"`
	Paused  bool           `json:"paused"`
	Mu      *sync.RWMutex
}

// Struct to be returned to a client
type Response struct {
	Message string `json:"message"`
}
