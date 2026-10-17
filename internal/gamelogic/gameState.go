package gamelogic

import (
	"sync"
)

type GameState struct {
	Players struct {
		Player1 Player `json:"player1"`
		Player2 Player `json:"player2"`
	} `json:"players"`
	Spectators   []Player `json:"spectators"`
	P1ActiveBots []*Bot   `json:"p1_active_bots"`
	P2ActiveBots []*Bot   `json:"p2_active_bots"`
	Turn         int      `json:"turn"`
	Mu           *sync.RWMutex
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
		Spectators:   make([]Player, 0),
		P1ActiveBots: make([]*Bot, 0, 2),
		P2ActiveBots: make([]*Bot, 0, 2),
		Turn:         0,
		Mu:           &sync.RWMutex{},
	}
}

// Initialize the game state before the game starts byt setting the active bots for both players and setting the turn counter to 1.
func (gs *GameState) BeginGame() {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	// Initialize the active bots for both players
	gs.P1ActiveBots = append(gs.P1ActiveBots, &gs.Players.Player1.Team[0], &gs.Players.Player1.Team[1])
	gs.P2ActiveBots = append(gs.P2ActiveBots, &gs.Players.Player2.Team[0], &gs.Players.Player2.Team[1])

	gs.Turn = 1 // Turn counter starts at 1
}

func (gs *GameState) GetPlayer1Name() string {
	return gs.Players.Player1.Username
}

func (gs *GameState) GetPlayer2Name() string {
	return gs.Players.Player2.Username
}

func (gs *GameState) GetPlayer1TeamNames() []string {
	names := make([]string, 0)
	for _, bot := range gs.Players.Player1.Team {
		names = append(names, bot.Name)
	}
	return names
}

func (gs *GameState) GetPlayer2TeamNames() []string {
	names := make([]string, 0)
	for _, bot := range gs.Players.Player2.Team {
		names = append(names, bot.Name)
	}
	return names
}
