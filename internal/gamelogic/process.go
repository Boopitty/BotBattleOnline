package gamelogic

import (
	"encoding/json"
	"fmt"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
)

// Request represents a command request from a player.
type Request struct {
	User    Player `json:"user"`
	Command string `json:"command"`
	Message string `json:"message"`
}

func Process(req []byte) []byte {
	// Parse the request string into a Request struct
	request := Request{}
	err := json.Unmarshal(req, &request)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command: %s", err)
		response := Response{Message: msg}
		return encoding.MakeJSONResponse(response)
	}

	response := commands(&request)
	return encoding.MakeJSONResponse(response)
}
