package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Boopitty/BotBattleOnline/internal/database"
	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/google/uuid"
)

// Struct to be returned to a client
type Response struct {
	Message string `json:"message"`
}

func commands(cfg *config, req *Request, w http.ResponseWriter, r *http.Request) any {
	// Handle the command based on the request
	switch req.Command {
	case "help":
		return Response{Message: "Available commands: help, profile, bots, attack, quit"}

	case "create-profile":
		dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      req.Login.Username,
			Password:  req.Login.Username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			encoding.RespondWithError(w, http.StatusBadRequest, err)
			return nil
		}
		return struct {
			username string
		}{
			username: dbUser.Name,
		}

	case "make-team":
		return Response{Message: "makeTeam not implemented yet"}

	case "new-game":
		gameState := gamelogic.NewGameState()
		gameState.Players[1] = req.User
		gameState.Players[2] = gamelogic.Player{Username: "BotPlayer", Team: make(map[int]gamelogic.Bot)}
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
