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
	case "help":
		return encoding.MakeJSONResponse(Response{Message: "Available commands: help, profile, bots, attack, quit"})

	case "select-bot":
		var bot *gamelogic.Bot
		if request.Input >= 0 {
			bot = gamelogic.MakeBot(gamelogic.BotNum(request.Input))
		}
		return encoding.MakeJSONResponse(Response{Message: fmt.Sprintf("Selected: %s", bot.Name)})

	case "make-team":
		return encoding.MakeJSONResponse(Response{Message: "makeTeam not implemented yet"})

	case "new-game":
		if _, exists := playerSession[conn]; exists {
			return encoding.MakeJSONResponse(Response{
				Message: "You are already in a game",
			})
		}

		// Check if the user is already in a game
		for id, gameState := range sessions {
			if gameState.Players.Player1.Conn == conn || gameState.Players.Player2.Conn == conn {
				playerSession[conn] = id
				return encoding.MakeJSONResponse(Response{
					Message: "You are already in a game",
				})
			}
		}

		// New unique session ID
		var sessionID uuid.UUID
		for {
			sessionID = uuid.New()
			_, exists := sessions[sessionID]
			if exists {
				continue
			} else {
				break
			}
		}

		// Store the gameState in a global variable and associate the connection with the session
		sessions[sessionID] = gamelogic.NewGameState()
		sessions[sessionID].Players.Player1 = request.User

		// Associate the connection with the session ID
		playerSession[conn] = sessionID

		fmt.Printf("New game started by player: %s\n", request.User.Username)
		fmt.Printf("New session ID: %s\n", sessionID)
		fmt.Printf("Current session count: %d\n", len(sessions))
		fmt.Printf("Current player count: %d\n", len(playerSession))

		return encoding.MakeJSONResponse(Response{
			Message: fmt.Sprintf("New game started by player: %s", request.User.Username),
		})

	case "join-game":
		// Check if the user is already in a game
		if _, exists := playerSession[conn]; exists {
			return encoding.MakeJSONResponse(Response{
				Message: "You are already in a game",
			})
		}

		// Find a game with less than 2 players and add the user to that game
		for id, gameState := range sessions {
			gameState.Mu.Lock()
			defer gameState.Mu.Unlock()

			if gameState.Players.Player1.Conn == nil {
				gameState.Players.Player1 = request.User
				playerSession[conn] = id

				fmt.Printf("Player 1: %s joined Player 2: %s\n", request.User.Username, gameState.Players.Player2.Username)
				fmt.Printf("Current player count: %d\n", len(playerSession))

				return encoding.MakeJSONResponse(Response{
					Message: fmt.Sprintf("You have joined a game against Player 2: %s", gameState.Players.Player2.Username),
				})

			} else if gameState.Players.Player2.Conn == nil {
				gameState.Players.Player2 = request.User
				playerSession[conn] = id

				fmt.Printf("Player 2: %s joined Player 1: %s\n", request.User.Username, gameState.Players.Player1.Username)
				fmt.Printf("Current player count: %d\n", len(playerSession))

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
		if sessionID, exists := playerSession[conn]; exists {
			delete(playerSession, conn)
			fmt.Printf("Current player count: %d\n", len(playerSession))

			if gameState, exists := sessions[sessionID]; exists {
				gameState.Mu.Lock()
				defer gameState.Mu.Unlock()

				if gameState.Players.Player1.Conn == conn {
					gameState.Players.Player1 = gamelogic.Player{}
				} else if gameState.Players.Player2.Conn == conn {
					gameState.Players.Player2 = gamelogic.Player{}
				}

				// If both players have left, remove the game state from sessions
				if gameState.Players.Player1.Conn == nil && gameState.Players.Player2.Conn == nil {
					delete(sessions, sessionID)
					fmt.Printf("Current session count: %d\n", len(sessions))

				}
			}
			return encoding.MakeJSONResponse(Response{Message: "You have left the game"})
		}
		return encoding.MakeJSONResponse(Response{Message: "You are not currently in a game"})

	case "pause":
		return encoding.MakeJSONResponse(Response{Message: "pause not implemented yet"})

	case "resume":
		return encoding.MakeJSONResponse(Response{Message: "resume not implemented yet"})

	case "profile":
		return encoding.MakeJSONResponse(Response{Message: "profile not implemented yet"})

	case "get-bot":
		bot := gamelogic.MakeBot(gamelogic.BotNum(request.Input))
		if bot == nil {
			return encoding.MakeJSONResponse(Response{
				Message: fmt.Sprintf("Unknown bot: %d", request.Input),
			})
		}
		return encoding.MakeJSONResponse(Response{
			Bot:     *bot,
			Message: fmt.Sprintf("Selected bot: %s", bot.Name),
		})

	case "attack":
		return encoding.MakeJSONResponse(Response{Message: "attack not implemented yet"})

	case "quit":
		return encoding.MakeJSONResponse(Response{Message: "Nice try, but you can't quit the game yet"})

	default:
		return encoding.MakeJSONResponse(Response{Message: fmt.Sprintf("Invalid command: %s", request.Command)})
	}
}
