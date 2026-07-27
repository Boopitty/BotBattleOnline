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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles a websocket request
func handleWS(cfg *config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // Upgrade the request into a websocket
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
			processed := Process(cfg, msg)
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
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
	}

	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Username:       req.Username,
		HashedPassword: hashedPass,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		encoding.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := struct {
		ID       uuid.UUID
		Username string
	}{
		ID:       user.ID,
		Username: user.Username,
	}

	encoding.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *config) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string
		Password string
	}{}

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		return
	}

	user, err := c.db.GetUser(r.Context(), req.Username)
	if err != nil {
		encoding.RespondWithError(w, http.StatusUnauthorized, err)
	}

	valid, err := auth.CheckPassHash(req.Password, user.HashedPassword)
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
	}
	if !valid {
		encoding.RespondWithError(w, http.StatusUnauthorized, nil)
		return
	}

	response := struct {
		Id       uuid.UUID
		Username string
	}{
		Id:       user.ID,
		Username: user.Username,
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

	user, err := c.db.GetUser(r.Context(), req.Username)
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

func (c *config) handleReset(w http.ResponseWriter, r *http.Request) {
	err := c.db.ResetUsers(r.Context())
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("!!! Database has been RESET !!!")
	encoding.RespondWithJSON(w, http.StatusOK, struct{}{})
}
