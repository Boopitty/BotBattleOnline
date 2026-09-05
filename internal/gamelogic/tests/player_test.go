package gamelogic_test

import (
	"sync"
	"testing"

	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
)

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		username string
		want     *gamelogic.Player
	}{
		{
			name:     "Default",
			username: "",
			want: &gamelogic.Player{
				Username: "",
				Team:     map[int]gamelogic.Bot{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gamelogic.NewPlayer(tt.username)

			if got.Username != tt.want.Username || len(got.Team) != len(tt.want.Team) {
				t.Errorf("NewPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGameState(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		want   *gamelogic.GameState
		wantMu *sync.RWMutex
	}{
		{
			name: "Default",
			want: &gamelogic.GameState{
				Players: struct {
					Player1 gamelogic.Player `json:"player1"`
					Player2 gamelogic.Player `json:"player2"`
				}{
					Player1: gamelogic.Player{},
					Player2: gamelogic.Player{},
				},
				Spectators: make([]gamelogic.Player, 0),
				ActiveBots: make([]gamelogic.Bot, 0, 4),
				Turn:       0,
				Mu:         &sync.RWMutex{},
			},
			wantMu: &sync.RWMutex{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gamelogic.NewGameState()

			if got.Players.Player1.Username != tt.want.Players.Player1.Username {
				t.Errorf("NewGameState().Players.Player1.Username = %v, want %v", got, tt.want)
			}
			if got.Players.Player2.Username != tt.want.Players.Player2.Username {
				t.Errorf("NewGameState().Players.Player2.Username = %v, want %v", got, tt.want)
			}
			if len(got.Players.Player1.Team) != len(tt.want.Players.Player1.Team) {
				t.Errorf("NewGameState().Players.Player1.Team = %v, want %v", got, tt.want)
			}
			if len(got.Players.Player2.Team) != len(tt.want.Players.Player2.Team) {
				t.Errorf("NewGameState().Players.Player2.Team = %v, want %v", got, tt.want)
			}
			if len(got.Spectators) != len(tt.want.Spectators) {
				t.Errorf("NewGameState().Spectators = %v, want %v", got, tt.want)
			}
			if len(got.ActiveBots) != len(tt.want.ActiveBots) {
				t.Errorf("NewGameState().ActiveBots = %v, want %v", got, tt.want)
			}
			if got.Turn != tt.want.Turn {
				t.Errorf("NewGameState().Turn = %v, want %v", got, tt.want)
			}
			if got.Mu == nil || tt.wantMu == nil {
				t.Errorf("NewGameState().Mu = %v, want %v", got, tt.want)
			}
		})
	}
}
