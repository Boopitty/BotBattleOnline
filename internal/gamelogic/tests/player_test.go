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
				Players: make(map[int]gamelogic.Player),
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
