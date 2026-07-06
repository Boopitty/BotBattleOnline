package main

import (
	"encoding/json"
	"fmt"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
)

func Process(cfg *config, req []byte) []byte {
	// Parse the request string into a Request struct
	request := Request{}
	err := json.Unmarshal(req, &request)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command: %s", err)
		response := Response{Message: msg}
		return encoding.MakeJSONResponse(response)
	}

	response := commands(cfg, &request)
	return encoding.MakeJSONResponse(response)
}
