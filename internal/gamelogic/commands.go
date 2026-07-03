package gamelogic

import (
	"fmt"
)

func commands(req *Request) string {
	// Handle the command based on the request
	switch req.Command {
	case "profile":
		return "profile not implemented yet"
	case "bots":
		return "bots not implemented yet"
	case "attack":
		return "attack not implemented yet"
	case "quit":
		return "Nice try, but you can't quit the game yet"
	default:
		return fmt.Sprintf("Invalid command: %s", req.Command)
	}
}
