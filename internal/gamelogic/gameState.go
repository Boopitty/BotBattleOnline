package gamelogic

import (
	"sync"
)

type GameState struct {
	Player1      Player   `json:"player1"`
	Player2      Player   `json:"player2"`
	Spectators   []Player `json:"spectators"`
	P1ActiveBots []*Bot   `json:"p1_active_bots"`
	P2ActiveBots []*Bot   `json:"p2_active_bots"`
	Turn         int      `json:"turn"`
	Mu           *sync.RWMutex
}

func NewGameState() *GameState {
	return &GameState{
		Player1:      Player{},
		Player2:      Player{},
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
	for _, bot := range gs.Player1.Team {
		if bot.Huge && (len(gs.P1ActiveBots) >= 1) {
			continue
		}
		gs.P1ActiveBots = append(gs.P1ActiveBots, &bot)
		if bot.Huge || (len(gs.P1ActiveBots) == 2) {
			break
		}
	}
	for _, bot := range gs.Player2.Team {
		if bot.Huge && (len(gs.P2ActiveBots) >= 1) {
			continue
		}
		gs.P2ActiveBots = append(gs.P2ActiveBots, &bot)
		if bot.Huge || (len(gs.P2ActiveBots) == 2) {
			break
		}
	}

	gs.Turn = 1 // Turn counter starts at 1
}

type playerNum int

const (
	player1 playerNum = 1
	player2 playerNum = 2
)

func (gs *GameState) GetPlayerName(p playerNum) string {
	switch p {
	case player1:
		return gs.Player1.Username
	case player2:
		return gs.Player2.Username
	default:
		return "None"
	}
}

func (gs *GameState) GetPlayerTeamNames(p playerNum) []string {
	names := make([]string, 0)
	switch p {
	case player1:
		for _, bot := range gs.Player1.Team {
			names = append(names, bot.Name)

		}
		return names

	case player2:
		for _, bot := range gs.Player2.Team {
			if bot.Huge && len(names) > 1 {
				continue
			}
			names = append(names, bot.Name)
			if bot.Huge {
				break
			}
		}
		return names

	default:
		return make([]string, 0)
	}
}
