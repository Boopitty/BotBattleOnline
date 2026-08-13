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

	response := commands(&request)
	return encoding.MakeJSONResponse(response)
}
