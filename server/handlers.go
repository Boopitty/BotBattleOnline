package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Boopitty/BotBattleOnline/internal/auth"
	"github.com/Boopitty/BotBattleOnline/internal/database"
	"github.com/Boopitty/BotBattleOnline/internal/encoding"
	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles a websocket request
func handleWS() func(http.ResponseWriter, *http.Request) {
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
				processed := Process(conn, []byte(`{"command":"leave-game"}`)) // Ensure the player leaves the game when the connection is closed
				log.Printf("Processed leaveGame response: %s\n", string(processed))
				break
			}
			log.Println("Received message:", string(msg))
			processed := Process(conn, msg)
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
		Log: "Echo Input: " + req.Command,
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

	token, err := auth.MakeJWT(user.ID, c.secret, time.Minute*5)
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}{
		Token:    token,
		Username: user.Username,
	}

	encoding.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *config) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	token, err := auth.MakeJWT(user.ID, c.secret, time.Hour*5)
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}{
		Token:    token,
		Username: user.Username,
	}
	encoding.RespondWithJSON(w, http.StatusOK, response)
}

func (c *config) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}{}

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		return
	}

	// Validate the token
	userID, err := auth.ValidateJWT(req.Token, c.secret)
	if err != nil {
		encoding.RespondWithError(w, http.StatusUnauthorized, err)
		return
	}

	// Get the user by ID
	user, err := c.db.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		encoding.RespondWithError(w, http.StatusUnauthorized, err)
		return
	}

	// Check if username and password match
	if user.Username != req.Username {
		encoding.RespondWithError(w, http.StatusUnauthorized, nil)
		return
	}
	match, err := auth.CheckPassHash(req.Password, user.HashedPassword)
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if !match {
		encoding.RespondWithError(w, http.StatusUnauthorized, nil)
		return
	}

	// Delete the user
	err = c.db.DeleteUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error Deleting user: %v", err)
		encoding.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	encoding.RespondWithJSON(w, http.StatusNoContent, struct{}{})
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

func (c *config) handleGetBot(w http.ResponseWriter, r *http.Request) {
	req := struct {
		BotName string `json:"BotName"`
	}{}

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		encoding.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Internal Error"))
		return
	}

	bot := gamelogic.MakeBot(gamelogic.BotType(req.BotName))
	if bot == nil {
		encoding.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unknown bot: %s", req.BotName))
		return
	}

	encoding.RespondWithJSON(w, http.StatusOK, bot)
}

func (c *config) handleSaveTeam(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Token string   `json:"Token"`
		Team  []string `json:"Team"`
	}{}

	success := encoding.DecodeJSON(w, r, &req)
	if !success {
		encoding.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Internal Error"))
		return
	}

	userID, err := auth.ValidateJWT(req.Token, c.secret)
	if err != nil {
		encoding.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
		return
	}

	err = c.db.DeleteUserTeam(r.Context(), userID) // Delete existing teams for the user
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Internal Error"))
		return
	}

	log.Printf("Saving team for user %s: \n%v", userID, req.Team)
	for i := 0; i < len(req.Team); i++ {
		_, err := c.db.CreateBot(r.Context(), database.CreateBotParams{
			ID:        uuid.New(),
			BotName:   req.Team[i],
			Skills:    make([]int32, 0), // Placeholder for skills, adjust as needed
			OwnerID:   userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			encoding.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}
		log.Printf("Bot saved: %s", req.Team[i])
	}

	botCount, err := c.db.GetNumBots(r.Context())
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Internal Error"))
		return
	}
	log.Printf("Total number of bots in the database: %d", botCount)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Team saved successfully",
	}
	encoding.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *config) handleGetTeamNames(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Internal Error"))
		return
	}

	userID, err := auth.ValidateJWT(token, c.secret)
	if err != nil {
		encoding.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
		return
	}

	team, err := c.db.GetUserTeamNames(r.Context(), userID)
	if err != nil {
		encoding.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Internal Error"))
		return
	}

	log.Printf("Retrieved team for user %s:\n %v", userID, team)

	response := struct {
		Team []string `json:"team"`
	}{
		Team: team,
	}
	encoding.RespondWithJSON(w, http.StatusOK, response)
}
