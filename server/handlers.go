package main

import (
	"encoding/json"
	"net/http"

	"github.com/Boopitty/BotBattleOnline/internal/pubsub"
)

func (c *config) handleCommand(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Command string `json:"command"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		pubsub.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response := struct {
		Log string `json:"log"`
	}{
		Log: "Command received: " + req.Command,
	}
	pubsub.RespondWithJSON(w, http.StatusOK, response)
}
