package main

import (
	"encoding/json"
	"fmt"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Incoming request struct.
// The Command field determines how to proceed with the request.
type Request struct {
	User      gamelogic.Player `json:"user"`
	TeamNames []string         `json:"teamnames"`
	Command   string           `json:"command"`
	Message   string           `json:"message"`
	Input     int              `json:"input"`
}

type OutPlayer struct {
	Name string   `json:"name"`
	Team []string `json:"team"`
}
type OutGameState struct {
	Player1 OutPlayer `json:"player1"`
	Player2 OutPlayer `json:"player2"`
}

// Struct to be returned to a client
type Response struct {
	Type      string       `json:"type"`
	Message   string       `json:"message"`
	GameState OutGameState `json:"gamestate"`
	Output    any          `json:"output"`
}

// Process websocket requests
func Process(conn *websocket.Conn, req []byte) {
	// Parse the request string into a Request struct
	var request Request
	err := json.Unmarshal(req, &request)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command: %s", err)
		response := Response{Message: msg}
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(response))
		return
	}
	request.User.Conn = conn // Placing the connection into the request struct makes copying player objects easier for later logic.

	switch request.Command {
	case "new-game":
		wsNewGame(conn, request)
		return

	case "join-game":
		wsJoinGame(conn, request)
		return

	case "leave-game":
		wsLeaveGame(conn)
		return

	case "get-teams":
		wsGetTeams(conn)
		return

	case "spectate":
		encoding.MakeJSONResponse(Response{Message: "spectate not implemented yet"})
		return

	case "leave-spectate":
		encoding.MakeJSONResponse(Response{Message: "leave-spectate not implemented yet"})
		return

	case "attack":
		encoding.MakeJSONResponse(Response{Message: "attack not implemented yet"})
		return

	default:
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
			Message: fmt.Sprintf("Invalid command: %s", request.Command),
		}))
		return
	}
}

func sendToAllPlayers(sessionID uuid.UUID, payload []byte) {
	if gameState, exists := activeGames.sessions[sessionID]; exists {
		gameState.Mu.RLock()
		defer gameState.Mu.RUnlock()

		if gameState.Player1.Conn != nil {
			gameState.Player1.Conn.WriteMessage(websocket.TextMessage, payload)
		}
		if gameState.Player2.Conn != nil {
			gameState.Player2.Conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}
