package gamelogic

import (
	"fmt"
)

type Response struct {
	Message string `json:"message"`
}

func commands(req *Request) Response {
	// Handle the command based on the request
	switch req.Command {
	case "help":
		return Response{Message: "Available commands: help, profile, bots, attack, quit"}

	case "make-team":
		return Response{Message: "makeTeam not implemented yet"}

	case "new-game":
		gameState := NewGameState()
		gameState.Players[1] = req.User
		gameState.Players[2] = Player{Username: "BotPlayer", Team: make(map[int]Bot)}
		// Here you would typically store the gameState in a global variable or a database
		// For example: globalGameState = gameState
		return Response{Message: "New game started with players: " + gameState.Players[1].Username + " and " + gameState.Players[2].Username}

	case "pause":
		return Response{Message: "pause not implemented yet"}

	case "resume":
		return Response{Message: "resume not implemented yet"}

	case "profile":
		return Response{Message: "profile not implemented yet"}

	case "bots":
		return Response{Message: "bots not implemented yet"}

	case "attack":
		return Response{Message: "attack not implemented yet"}

	case "quit":
		return Response{Message: "Nice try, but you can't quit the game yet"}

	default:
		return Response{Message: fmt.Sprintf("Invalid command: %s", req.Command)}
	}
}
