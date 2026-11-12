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

	// Find a game with an empty player slot and add the user to that game
	for id, gameState := range activeGames.sessions {
		success := gameState.JoinGame(request.User)
		if !success {
			continue
		}
		activePlayers.sessions[conn] = id // Associate the connection with the session ID

		sendToAll(id, encoding.MakeJSONResponse(Response{
			Message: fmt.Sprintf("%s has joined the game", request.User.Username),
		}))

		// If both players are now present, begin the game
		if gameState.Player1.Conn != nil && gameState.Player2.Conn != nil {
			gameState.BeginGame()
			sendGSToAll(id)
		}
		break
	}
	conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "No available games to join"}))
}

func wsLeaveGame(conn *websocket.Conn) {
	if sessionID, exists := activePlayers.sessions[conn]; exists {
		delete(activePlayers.sessions, conn)

		if gameState, exists := activeGames.sessions[sessionID]; exists {
			gameState.LeaveGame(conn) // Remove the player from the game state
			sendToAll(sessionID, encoding.MakeJSONResponse(Response{
				Message: "A player has left the game",
			}))

			// If both players have left, remove the game state from sessions
			if gameState.Player1.Conn == nil && gameState.Player2.Conn == nil {
				delete(activeGames.sessions, sessionID)
			} else {
				sendGSToAll(sessionID)
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

func sendToAll(sessionID uuid.UUID, payload []byte) {
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

func sendGSToAll(sessionID uuid.UUID) {
	if gameState, exists := activeGames.sessions[sessionID]; exists {
		gameState.Mu.RLock()
		defer gameState.Mu.RUnlock()

		response := Response{
			Type:    "gamestate",
			Message: "Game State Update",
			GameState: OutGameState{
				Player1: OutPlayer{
					Name: gameState.GetPlayerName(1),
					Team: gameState.GetPlayerTeamNames(1),
				},
				Player2: OutPlayer{
					Name: gameState.GetPlayerName(2),
					Team: gameState.GetPlayerTeamNames(2),
				},
				P1Actives: gameState.P1ActiveBots,
				P2Actives: gameState.P2ActiveBots,
				Finished:  gameState.IsGameOver(),
			},
		}
		payload := encoding.MakeJSONResponse(response)

		if gameState.Player1.Conn != nil {
			gameState.Player1.Conn.WriteMessage(websocket.TextMessage, payload)
		}
		if gameState.Player2.Conn != nil {
			gameState.Player2.Conn.WriteMessage(websocket.TextMessage, payload)
		}
		for _, spectatorConn := range gameState.Spectators {
			spectatorConn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}
