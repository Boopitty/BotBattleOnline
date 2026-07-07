package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
)

// Incoming reques struct
type Request struct {
	User    gamelogic.Player `json:"user"`
	Command string           `json:"command"`
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
