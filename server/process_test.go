package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/gorilla/websocket"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req  Request
		want Response
	}{
		{
			name: "Unknown command",
			req: Request{
				Command: "test-command",
			},
			want: Response{
				Message: "Invalid command: test-command",
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
			want: Response{
				Message: "Invalid Team",
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
				TeamNames: []string{"Assault"},
				Command:   "new-game",
			},
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
				TeamNames: []string{"Assault"},
				Command:   "join-game",
			},
			want: Response{
				Type:    "gamestate",
				Message: "Sending Gamestate",
				GameState: OutGameState{
					Player1: OutPlayer{
						Name: "testuser1",
						Team: []string{"Assault"},
					},
					Player2: OutPlayer{
						Name: "testuser2",
						Team: []string{"Assault"},
					},
				},
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

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					t.Errorf("Failed to upgrade websocket: %v", err)
					return
				}
				defer conn.Close()
				Process(conn, reqBytes)
			}))
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Fatalf("Failed to connect to websocket: %v", err)
			}
			defer client.Close()

			// check the written message
			_, reader, err := client.NextReader()
			if err != nil {
				t.Fatalf("Error in conn.NextReader: %v", err)
			}

			got, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("Error in io.ReadAll: %v", err)
			}

			var gotResp Response
			if err := json.Unmarshal(got, &gotResp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
			if gotResp.Message != tt.want.Message {
				t.Errorf("Process() = %v, want %v", gotResp, tt.want)
			}
		})
	}
}
