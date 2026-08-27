package main

import (
	"encoding/json"
	"fmt"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
)

// Incoming request struct.
// The Command field determines how to proceed with the request.
type Request struct {
	User    gamelogic.Player `json:"user"`
	Command string           `json:"command"`
	Message string           `json:"message"`
}

// Struct to be returned to a client
type Response struct {
	Message string `json:"message"`
}

// Process websocket requests
func Process(cfg *config, req []byte) []byte {
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
		return encoding.MakeJSONResponse(Response{Message: "selectBot not implemented yet"})

	case "make-team":
		return encoding.MakeJSONResponse(Response{Message: "makeTeam not implemented yet"})

	case "new-game":
		gameState := gamelogic.NewGameState()
		gameState.Players[1] = request.User
		gameState.Players[2] = gamelogic.Player{Username: "BotPlayer", Team: make(map[int]gamelogic.Bot)}
		// Here you would typically store the gameState in a global variable or a database
		// For example: globalGameState = gameState
		return encoding.MakeJSONResponse(Response{
			Message: "New game started with players: " + gameState.Players[1].Username + " and " + gameState.Players[2].Username,
		})

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
