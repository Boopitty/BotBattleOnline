package main

import (
	"fmt"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func wsNewGame(conn *websocket.Conn, request Request) {
	if len(request.TeamNames) < 1 {
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
			Message: "Invalid Team",
		}))
		return
	}

	activePlayers.mu.Lock()
	defer activePlayers.mu.Unlock()

	// Check if the user's connection is already associated with a game session
	if _, exists := activePlayers.sessions[conn]; exists {
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
			Message: "You are already in a game",
		}))
		return
	}

	activeGames.mu.Lock()
	defer activeGames.mu.Unlock()

	// Check if the user's connection is already associated with a game session in activeGames
	for id, gameState := range activeGames.sessions {
		if gameState.Player1.Conn == conn || gameState.Player2.Conn == conn {
			activePlayers.sessions[conn] = id // Re-associate the connection with the existing session ID
			conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
				Message: "You are already in a game",
			}))
			return
		}
	}

	// Create new unique session ID
	var sessionID uuid.UUID
	for { // Ensure there are no duplicates
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
	activeGames.sessions[sessionID].Player1 = request.User
	for _, botName := range request.TeamNames {
		bot := gamelogic.MakeBot(gamelogic.BotType(botName))
		if bot != nil {
			activeGames.sessions[sessionID].Player1.Team = append(activeGames.sessions[sessionID].Player1.Team, *bot)
		}
	}

	// Associate the connection with the session ID
	activePlayers.sessions[conn] = sessionID

	fmt.Printf("New game started by player: %s\n", request.User.Username)

	conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
		Message: fmt.Sprintf("New game started by player: %s", request.User.Username),
	}))
}

func wsJoinGame(conn *websocket.Conn, request Request) {
	if len(request.TeamNames) < 1 {
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
			Message: "Invalid Team",
		}))
		return
	}

	activePlayers.mu.Lock()
	defer activePlayers.mu.Unlock()

	// Check if the user's connection is already associated with a game session
	if _, exists := activePlayers.sessions[conn]; exists {
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
			Message: "You are already in a game",
		}))
		return
	}

	activeGames.mu.Lock()
	defer activeGames.mu.Unlock()

	// Check if the user's connection is already associated with a game session in activeGames
	for id, gameState := range activeGames.sessions {
		if gameState.Player1.Conn == conn || gameState.Player2.Conn == conn {
			activePlayers.sessions[conn] = id // Re-associate the connection with the existing session ID
			conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
				Message: "You are already in a game",
			}))
			return
		}
	}

	found := false
	// Find a game with an empty player slot and add the user to that game
	for !found {
		for id, gameState := range activeGames.sessions {
			if gameState.Player1.Conn == nil {
				gameState.Player1 = request.User // Add the user as Player 1
				found = true
				for _, botName := range request.TeamNames { // Add the user's team to Player 1's team
					bot := gamelogic.MakeBot(gamelogic.BotType(botName))
					if bot != nil {
						gameState.Player1.Team = append(gameState.Player1.Team, *bot)
					}
				}

				activePlayers.sessions[conn] = id // Associate the connection with the session ID

				fmt.Printf("Player 1: %s joined Player 2: %s\n", request.User.Username, gameState.Player2.Username)
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

				conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
					Message: fmt.Sprintf("You have joined a game against Player 2: %s", gameState.Player2.Username),
				}))
				return

			} else if gameState.Player2.Conn == nil {
				gameState.Player2 = request.User // Add the user as Player 2
				found = true
				for _, botName := range request.TeamNames { // Add the user's team to Player 2's team
					bot := gamelogic.MakeBot(gamelogic.BotType(botName))
					if bot != nil {
						gameState.Player2.Team = append(gameState.Player2.Team, *bot)
					}
				}

				activePlayers.sessions[conn] = id // Associate the connection with the session ID

				fmt.Printf("Player 2: %s joined Player 1: %s\n", request.User.Username, gameState.Player1.Username)
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

				conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
					Message: fmt.Sprintf("You have joined a game against Player 1: %s", gameState.Player1.Username),
				}))
				return
			}
		}
	}
	conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "No available games to join"}))
}

func wsLeaveGame(conn *websocket.Conn) {
	if sessionID, exists := activePlayers.sessions[conn]; exists {
		delete(activePlayers.sessions, conn)

		if gameState, exists := activeGames.sessions[sessionID]; exists {
			gameState.Mu.Lock()
			defer gameState.Mu.Unlock()

			if gameState.Player1.Conn == conn {
				gameState.Player1 = gamelogic.Player{}
			} else if gameState.Player2.Conn == conn {
				gameState.Player2 = gamelogic.Player{}
			}

			// If both players have left, remove the game state from sessions
			if gameState.Player1.Conn == nil && gameState.Player2.Conn == nil {
				delete(activeGames.sessions, sessionID)
			}
		}
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "You have left the game"}))
	}
	conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "You are not currently in a game"}))
}

func wsGetTeams(conn *websocket.Conn) {
	if sessionID, exists := activePlayers.sessions[conn]; exists {
		if gameState, exists := activeGames.sessions[sessionID]; exists {
			gameState.Mu.RLock()
			defer gameState.Mu.RUnlock()

			// Return the teams of both players
			conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
				Message: "Teams retrieved successfully",
				Output: map[string]gamelogic.Player{
					"player1": gameState.Player1,
					"player2": gameState.Player2,
				},
			}))
			return
		}
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "Game state not found"}))
	}
	conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "You are not currently in a game"}))
}
