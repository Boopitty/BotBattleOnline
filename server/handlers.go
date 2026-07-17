package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Boopitty/BotBattleOnline/internal/auth"
	"github.com/Boopitty/BotBattleOnline/internal/database"
	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
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

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		return
	}

	response := struct {
		Log string `json:"log"`
	}{
		Log: "Command received: " + req.Command,
	}
	encoding.RespondWithJSON(w, http.StatusOK, response)
}

func (c *config) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		return
	}

	hashedPass, err := auth.HashPass(req.Password)
	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      req.Username,
		Password:  hashedPass,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		encoding.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := struct {
		ID   uuid.UUID
		Name string
	}{
		ID:   user.ID,
		Name: user.Name,
	}

	encoding.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *config) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		return
	}

	user, err := c.db.LoginUser(r.Context(), database.LoginUserParams{
		Name:     req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Error finding user: %v", err)
		encoding.RespondWithError(w, http.StatusUnauthorized, err)
		return
	}
	err = c.db.DeleteUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error Deleting user: %v", err)
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
