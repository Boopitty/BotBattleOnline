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

	var sessionID uuid.UUID
	for { // Ensure there are no duplicates for the session ID
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
	activeGames.sessions[sessionID].JoinGame(request.User)

	// Associate the connection with the session ID
	activePlayers.sessions[conn] = sessionID

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

	// Check if a game session already has this connection associated with it
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

	if _, exists := activePlayers.sessions[conn]; !exists {
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{
			Message: "No available games to join",
		}))
	}
}

func wsLeaveGame(conn *websocket.Conn) {
	if sessionID, exists := activePlayers.sessions[conn]; exists {
		delete(activePlayers.sessions, conn)

		if gameState, exists := activeGames.sessions[sessionID]; exists {
			gameState.LeaveGame(conn)
			sendToAll(sessionID, encoding.MakeJSONResponse(Response{
				Message: "A player has left the game",
			}))

			// If both players have left, remove the game state from sessions
			if gameState.Player1.Conn == nil && gameState.Player2.Conn == nil {
				sendGSToAll(sessionID)
				delete(activeGames.sessions, sessionID)
			} else {
				sendGSToAll(sessionID)
			}
		}
		conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "You have left the game"}))
		return
	}
	conn.WriteMessage(websocket.TextMessage, encoding.MakeJSONResponse(Response{Message: "You are not currently in a game"}))
}

// Send payload to all players in a session
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

// Send game state to all players and spectators in a session
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
				P1Idle:    gameState.GetIdleNames(1),
				P2Idle:    gameState.GetIdleNames(2),
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
