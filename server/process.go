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
	User    gamelogic.Player `json:"user"`
	Command string           `json:"command"`
	Message string           `json:"message"`
	Input   int              `json:"input"`
}

// Struct to be returned to a client
type Response struct {
	Bot     gamelogic.Bot `json:"bot,omitempty"`
	Message string        `json:"message"`
	Output  any           `json:"output"`
}

// Process websocket requests
func Process(conn *websocket.Conn, req []byte) []byte {
	// Parse the request string into a Request struct
	var request Request
	err := json.Unmarshal(req, &request)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command: %s", err)
		response := Response{Message: msg}
		return encoding.MakeJSONResponse(response)
	}
	request.User.Conn = conn // Associate the connection with the user to be stored in the game state

	switch request.Command {
	case "new-game":
		activePlayers.mu.Lock()
		defer activePlayers.mu.Unlock()

		// Check if the user's connection is already associated with a game session
		if _, exists := activePlayers.sessions[conn]; exists {
			return encoding.MakeJSONResponse(Response{
				Message: "You are already in a game",
			})
		}

		// Check if the user's connection is already associated with a game session in activeGames
		for id, gameState := range activeGames.sessions {
			if gameState.Players.Player1.Conn == conn || gameState.Players.Player2.Conn == conn {
				activePlayers.sessions[conn] = id // Re-associate the connection with the existing session ID
				return encoding.MakeJSONResponse(Response{
					Message: "You are already in a game",
				})
			}
		}

		activeGames.mu.Lock()
		defer activeGames.mu.Unlock()

		// Create new unique session ID
		var sessionID uuid.UUID
		for {
			sessionID = uuid.New()
			_, exists := activeGames.sessions[sessionID]
			if exists {
				continue
			} else {
				break
			}
		}

		// Store the gameState in a global variable and associate the connection with the session
		activeGames.sessions[sessionID] = gamelogic.NewGameState()
		activeGames.sessions[sessionID].Players.Player1 = request.User

		// Associate the connection with the session ID
		activePlayers.sessions[conn] = sessionID

		fmt.Printf("New game started by player: %s\n", request.User.Username)

		return encoding.MakeJSONResponse(Response{
			Message: fmt.Sprintf("New game started by player: %s", request.User.Username),
		})

	case "join-game":
		activePlayers.mu.Lock()
		defer activePlayers.mu.Unlock()

		// Check if the user's connection is already associated with a game session
		if _, exists := activePlayers.sessions[conn]; exists {
			return encoding.MakeJSONResponse(Response{
				Message: "You are already in a game",
			})
		}

		activeGames.mu.Lock()
		defer activeGames.mu.Unlock()

		// Check if the user's connection is already associated with a game session in activeGames
		for id, gameState := range activeGames.sessions {
			if gameState.Players.Player1.Conn == conn || gameState.Players.Player2.Conn == conn {
				activePlayers.sessions[conn] = id // Re-associate the connection with the existing session ID
				return encoding.MakeJSONResponse(Response{
					Message: "You are already in a game",
				})
			}
		}

		// Find a game with less than 2 players and add the user to that game
		for id, gameState := range activeGames.sessions {
			if gameState.Players.Player1.Conn == nil {
				gameState.Players.Player1 = request.User
				activePlayers.sessions[conn] = id

				fmt.Printf("Player 1: %s joined Player 2: %s\n", request.User.Username, gameState.Players.Player2.Username)

				return encoding.MakeJSONResponse(Response{
					Message: fmt.Sprintf("You have joined a game against Player 2: %s", gameState.Players.Player2.Username),
				})

			} else if gameState.Players.Player2.Conn == nil {
				gameState.Players.Player2 = request.User
				activePlayers.sessions[conn] = id

				fmt.Printf("Player 2: %s joined Player 1: %s\n", request.User.Username, gameState.Players.Player1.Username)

				return encoding.MakeJSONResponse(Response{
					Message: fmt.Sprintf("You have joined a game against Player 1: %s", gameState.Players.Player1.Username),
				})
			}
		}
		return encoding.MakeJSONResponse(Response{Message: "No available games to join"})

	case "spectate":
		return encoding.MakeJSONResponse(Response{Message: "spectate not implemented yet"})

	case "leave-spectate":
		return encoding.MakeJSONResponse(Response{Message: "leave-spectate not implemented yet"})

	case "leave-game":
		if sessionID, exists := activePlayers.sessions[conn]; exists {
			delete(activePlayers.sessions, conn)

			if gameState, exists := activeGames.sessions[sessionID]; exists {
				gameState.Mu.Lock()
				defer gameState.Mu.Unlock()

				if gameState.Players.Player1.Conn == conn {
					gameState.Players.Player1 = gamelogic.Player{}
				} else if gameState.Players.Player2.Conn == conn {
					gameState.Players.Player2 = gamelogic.Player{}
				}

				// If both players have left, remove the game state from sessions
				if gameState.Players.Player1.Conn == nil && gameState.Players.Player2.Conn == nil {
					delete(activeGames.sessions, sessionID)
				}
			}
			return encoding.MakeJSONResponse(Response{Message: "You have left the game"})
		}
		return encoding.MakeJSONResponse(Response{Message: "You are not currently in a game"})

	case "pause":
		return encoding.MakeJSONResponse(Response{Message: "pause not implemented yet"})

	case "resume":
		return encoding.MakeJSONResponse(Response{Message: "resume not implemented yet"})

	case "attack":
		return encoding.MakeJSONResponse(Response{Message: "attack not implemented yet"})

	default:
		return encoding.MakeJSONResponse(Response{Message: fmt.Sprintf("Invalid command: %s", request.Command)})
	}
}
