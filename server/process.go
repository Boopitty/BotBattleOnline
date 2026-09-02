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
	Message string `json:"message"`
	Output  any    `json:"output"`
}

// Process websocket requests
func Process(cfg *config, req []byte, conn *websocket.Conn) []byte {
	// Parse the request string into a Request struct
	request := Request{}
	err := json.Unmarshal(req, &request)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command: %s", err)
		response := Response{Message: msg}
		return encoding.MakeJSONResponse(response)
	}

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
		// Check if the user is already in a game
		for id, gameState := range sessions {
			if _, exists := gameState.Players[conn]; exists {
				return encoding.MakeJSONResponse(Response{
					Message: "You are already in a game",
					Output:  id,
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

		// Store the gameState in a global variable
		sessions[sessionID] = gamelogic.NewGameState()
		sessions[sessionID].Players[conn] = request.User
		return encoding.MakeJSONResponse(Response{
			Message: "New game started with player: " + sessions[sessionID].Players[conn].Username,
			Output:  sessionID,
		})

	case "join-game":
		return encoding.MakeJSONResponse(Response{Message: "joinGame not implemented yet"})

	case "spectate-game":
		return encoding.MakeJSONResponse(Response{Message: "spectateGame not implemented yet"})

	case "leave-game":
		return encoding.MakeJSONResponse(Response{Message: "leaveGame not implemented yet"})

	case "pause":
		return encoding.MakeJSONResponse(Response{Message: "pause not implemented yet"})

	case "resume":
		return encoding.MakeJSONResponse(Response{Message: "resume not implemented yet"})

	case "profile":
		return encoding.MakeJSONResponse(Response{Message: "profile not implemented yet"})

	case "bots":
		return encoding.MakeJSONResponse(Response{Message: "bots not implemented yet"})

	case "attack":
		return encoding.MakeJSONResponse(Response{Message: "attack not implemented yet"})

	case "quit":
		return encoding.MakeJSONResponse(Response{Message: "Nice try, but you can't quit the game yet"})

	default:
		return encoding.MakeJSONResponse(Response{Message: fmt.Sprintf("Invalid command: %s", request.Command)})
	}
}
