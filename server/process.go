package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
)

// Incoming request struct.
// When request is decoded, it fills as much of this struct as possible.
// The Command field determines how to proceed with the request, while using only the relevant fields.
type Request struct {
	User    gamelogic.Player `json:"user"`    // manditory
	Command string           `json:"command"` // manditory
	Message string           `json:"message"`
	Login   struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
}

func Process(cfg *config, req []byte, w http.ResponseWriter, r *http.Request) []byte {
	// Parse the request string into a Request struct
	request := Request{}
	err := json.Unmarshal(req, &request)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command: %s", err)
		response := Response{Message: msg}
		return encoding.MakeJSONResponse(response)
	}

	response := commands(cfg, &request, w, r)
	return encoding.MakeJSONResponse(response)
}
