package main

import (
	"encoding/json"
	"testing"

	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func TestProcess(t *testing.T) {
	playerSession = make(map[*websocket.Conn]uuid.UUID) // Map of player connections to session IDs
	sessions = make(map[uuid.UUID]*gamelogic.GameState) // List of gamestates

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req  Request
		conn *websocket.Conn
		want Response
	}{
		{
			name: "Unknown command",
			req: Request{
				Command: "test-command",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "Invalid command: test-command",
			},
		},
		{
			name: "Help command",
			req: Request{
				Command: "help",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "Available commands: help, profile, bots, attack, quit",
			},
		},
		{
			name: "join-game with no existing game",
			req: Request{
				User: gamelogic.Player{
					Username: "testFailJoin",
				},
				Command: "join-game",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "No available games to join",
			},
		},
		{
			name: "leave-game while not in a game",
			req: Request{
				User: gamelogic.Player{
					Username: "testFailLeave",
				},
				Command: "leave-game",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "You are not currently in a game",
			},
		},
		{
			name: "new-game command",
			req: Request{
				User: gamelogic.Player{
					Username: "testuser1",
				},
				Command: "new-game",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "New game started by player: testuser1",
			},
		},
		{
			name: "join-game command",
			req: Request{
				User: gamelogic.Player{
					Username: "testuser2",
				},
				Command: "join-game",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "You have joined a game against Player 1: testuser1",
			},
		},
		{
			name: "leave-game command",
			req: Request{
				User: gamelogic.Player{
					Username: "testuser1",
				},
				Command: "leave-game",
			},
			conn: &websocket.Conn{},
			want: Response{
				Message: "You are not currently in a game",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBytes, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			got := Process(tt.conn, reqBytes)

			var gotResp Response
			if err := json.Unmarshal(got, &gotResp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if gotResp != tt.want {
				t.Errorf("Process() = %v, want %v", gotResp, tt.want)
			}
		})
	}
}
