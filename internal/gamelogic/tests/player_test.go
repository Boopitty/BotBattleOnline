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
				Team:     make([]gamelogic.Bot, 0),
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
				Player1:      gamelogic.Player{},
				Player2:      gamelogic.Player{},
				Spectators:   make([]gamelogic.Player, 0),
				P1ActiveBots: make([]*gamelogic.Bot, 0, 2),
				P2ActiveBots: make([]*gamelogic.Bot, 0, 2),
				Turn:         0,
				Mu:           &sync.RWMutex{},
			},
			wantMu: &sync.RWMutex{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gamelogic.NewGameState()

			if got.Player1.Username != tt.want.Player1.Username {
				t.Errorf("NewGameState().Players.Player1.Username = %v, want %v", got, tt.want)
			}
			if got.Player2.Username != tt.want.Player2.Username {
				t.Errorf("NewGameState().Players.Player2.Username = %v, want %v", got, tt.want)
			}
			if len(got.Player1.Team) != len(tt.want.Player1.Team) {
				t.Errorf("NewGameState().Players.Player1.Team = %v, want %v", got, tt.want)
			}
			if len(got.Player2.Team) != len(tt.want.Player2.Team) {
				t.Errorf("NewGameState().Players.Player2.Team = %v, want %v", got, tt.want)
			}
			if len(got.Spectators) != len(tt.want.Spectators) {
				t.Errorf("NewGameState().Spectators = %v, want %v", got, tt.want)
			}
			if len(got.P1ActiveBots) != len(tt.want.P1ActiveBots) {
				t.Errorf("NewGameState().P1ActiveBots = %v, want %v", got, tt.want)
			}
			if len(got.P2ActiveBots) != len(tt.want.P2ActiveBots) {
				t.Errorf("NewGameState().P2ActiveBots = %v, want %v", got, tt.want)
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
