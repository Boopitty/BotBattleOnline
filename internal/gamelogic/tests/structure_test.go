package gamelogic_test

import (
	"sync"
	"testing"

	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
)

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want *gamelogic.Player
	}{
		{
			name: "NewPlayer returns a Player with empty username and empty team",
			want: &gamelogic.Player{
				Username: "",
				Team:     map[int]gamelogic.Bot{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gamelogic.NewPlayer()

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
			name: "NewGameState returns a GameState with empty players and paused set to false",
			want: &gamelogic.GameState{
				Players: []gamelogic.Player{},
				Paused:  false,
				Mu:      &sync.RWMutex{},
			},
			wantMu: &sync.RWMutex{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gamelogic.NewGameState()

			if len(got.Players) != len(tt.want.Players) || got.Paused != tt.want.Paused {
				t.Errorf("NewGameState() = %v, want %v", got, tt.want)
			}
		})
	}
}
