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
		if len(request.TeamNames) < 1 {
			return encoding.MakeJSONResponse(Response{
				Message: "Invalid Team",
			})
		}

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

		// Set the user as Player 1 in the game state and add their team
		activeGames.sessions[sessionID].Players.Player1 = request.User
		for _, botName := range request.TeamNames {
			bot := gamelogic.MakeBot(gamelogic.BotType(botName))
			if bot != nil {
				activeGames.sessions[sessionID].Players.Player1.Team = append(activeGames.sessions[sessionID].Players.Player1.Team, *bot)
			}
		}

		// Associate the connection with the session ID
		activePlayers.sessions[conn] = sessionID

		fmt.Printf("New game started by player: %s\n", request.User.Username)

		return encoding.MakeJSONResponse(Response{
			Message: fmt.Sprintf("New game started by player: %s", request.User.Username),
		})

	case "join-game":
		if len(request.TeamNames) < 1 {
			return encoding.MakeJSONResponse(Response{
				Message: "Invalid Team",
			})
		}

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

		found := false
		// Find a game with an empty player slot and add the user to that game
		for !found {
			for id, gameState := range activeGames.sessions {
				if gameState.Players.Player1.Conn == nil {
					gameState.Players.Player1 = request.User // Add the user as Player 1
					found = true
					for _, botName := range request.TeamNames { // Add the user's team to Player 1's team
						bot := gamelogic.MakeBot(gamelogic.BotType(botName))
						if bot != nil {
							gameState.Players.Player1.Team = append(gameState.Players.Player1.Team, *bot)
						}
					}

					activePlayers.sessions[conn] = id // Associate the connection with the session ID

					fmt.Printf("Player 1: %s joined Player 2: %s\n", request.User.Username, gameState.Players.Player2.Username)
					gameState.BeginGame()
					sendToAllPlayers(id, encoding.MakeJSONResponse(Response{
						Type:    "gamestate",
						Message: "Sending Gamestate",
						GameState: OutGameState{
							Player1: OutPlayer{
								Name: gameState.GetPlayer1Name(),
								Team: gameState.GetPlayer1TeamNames(),
							},
							Player2: OutPlayer{
								Name: gameState.GetPlayer2Name(),
								Team: gameState.GetPlayer2TeamNames(),
							},
						},
					}))

					return encoding.MakeJSONResponse(Response{
						Message: fmt.Sprintf("You have joined a game against Player 2: %s", gameState.Players.Player2.Username),
					})

				} else if gameState.Players.Player2.Conn == nil {
					gameState.Players.Player2 = request.User // Add the user as Player 2
					found = true
					for _, botName := range request.TeamNames { // Add the user's team to Player 2's team
						bot := gamelogic.MakeBot(gamelogic.BotType(botName))
						if bot != nil {
							gameState.Players.Player2.Team = append(gameState.Players.Player2.Team, *bot)
						}
					}

					activePlayers.sessions[conn] = id // Associate the connection with the session ID

					fmt.Printf("Player 2: %s joined Player 1: %s\n", request.User.Username, gameState.Players.Player1.Username)
					gameState.BeginGame()
					sendToAllPlayers(id, encoding.MakeJSONResponse(Response{
						Type:    "gamestate",
						Message: "Sending Gamestate",
						GameState: OutGameState{
							Player1: OutPlayer{
								Name: gameState.GetPlayer1Name(),
								Team: gameState.GetPlayer1TeamNames(),
							},
							Player2: OutPlayer{
								Name: gameState.GetPlayer2Name(),
								Team: gameState.GetPlayer2TeamNames(),
							},
						},
					}))

					return encoding.MakeJSONResponse(Response{
						Message: fmt.Sprintf("You have joined a game against Player 1: %s", gameState.Players.Player1.Username),
					})
				}
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

	case "get-teams":
		if sessionID, exists := activePlayers.sessions[conn]; exists {
			if gameState, exists := activeGames.sessions[sessionID]; exists {
				gameState.Mu.RLock()
				defer gameState.Mu.RUnlock()

				// Return the teams of both players
				return encoding.MakeJSONResponse(Response{
					Message: "Teams retrieved successfully",
					Output: map[string]gamelogic.Player{
						"player1": gameState.Players.Player1,
						"player2": gameState.Players.Player2,
					},
				})
			}
			return encoding.MakeJSONResponse(Response{Message: "Game state not found"})
		}
		return encoding.MakeJSONResponse(Response{Message: "You are not currently in a game"})

	case "initialize-game":
		// Initialize the game state for the current session
		if sessionID, exists := activePlayers.sessions[conn]; exists {
			if gameState, exists := activeGames.sessions[sessionID]; exists {
				// Check if both players have set their teams
				if gameState.Players.Player1.Team == nil || gameState.Players.Player2.Team == nil {
					return encoding.MakeJSONResponse(Response{Message: "Waiting for both players to set their teams"})
				}

				// Initialize the game state
				gameState.BeginGame()
				return encoding.MakeJSONResponse(Response{Message: "Game initialized successfully"})
			}
			return encoding.MakeJSONResponse(Response{Message: "Game state not found"})
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

func sendToAllPlayers(sessionID uuid.UUID, payload []byte) {
	if gameState, exists := activeGames.sessions[sessionID]; exists {
		gameState.Mu.RLock()
		defer gameState.Mu.RUnlock()

		if gameState.Players.Player1.Conn != nil {
			gameState.Players.Player1.Conn.WriteMessage(websocket.TextMessage, payload)
		}
		if gameState.Players.Player2.Conn != nil {
			gameState.Players.Player2.Conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}
