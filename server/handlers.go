package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/gorilla/websocket"
)

func handleWS(cfg *config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading message: %v", err)
				break
			}
			log.Println("Received message:", string(msg))
			processed := Process(cfg, msg, w, r)
			conn.WriteMessage(websocket.TextMessage, processed)
		}
	}
}

func (c *config) handleCommand(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Command string `json:"command"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		encoding.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := struct {
		Log string `json:"log"`
	}{
		Log: "Command received: " + req.Command,
	}
	encoding.RespondWithJSON(w, http.StatusOK, response)
}
