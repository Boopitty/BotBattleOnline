package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var activeGames = struct {
	sessions map[uuid.UUID]*gamelogic.GameState
	mu       *sync.Mutex
}{
	sessions: make(map[uuid.UUID]*gamelogic.GameState),
	mu:       &sync.Mutex{},
}

var activePlayers = struct {
	sessions map[*websocket.Conn]uuid.UUID
	mu       *sync.Mutex
}{
	sessions: make(map[*websocket.Conn]uuid.UUID),
	mu:       &sync.Mutex{},
}

func main() {
	cfg := createConfig()

	mux := http.NewServeMux()                                          // Mutex for the server
	appHandler := http.FileServer(http.Dir(cfg.filepathRoot))          // Serve the static files from the specified root directory
	statisticsHandler := http.FileServer(http.Dir(cfg.statisticsRoot)) // Serve the static files from the specified root directory

	// Handle the front page
	mux.Handle("/", http.StripPrefix("/", appHandler))
	mux.Handle("/statistics/", http.StripPrefix("/statistics/", statisticsHandler))

	// Handle the API endpoints
	mux.HandleFunc("/api/command", cfg.handleCommand)
	mux.HandleFunc("POST /api/createUser", cfg.handleCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)
	mux.HandleFunc("DELETE /api/deleteUser", cfg.handleDeleteUser)
	mux.HandleFunc("DELETE /api/resetUsers", cfg.handleReset)
	mux.HandleFunc("GET /api/getBot", cfg.handleGetBot)
	mux.HandleFunc("POST /api/saveTeam", cfg.handleSaveTeam)
	mux.HandleFunc("GET /api/getTeamNames", cfg.handleGetTeamNames)
	mux.HandleFunc("/ws", handleWS())

	// Handle Admin endpoints
	mux.HandleFunc("GET /admin/getNumUsers", cfg.handleGetNumUsers)
	mux.HandleFunc("GET /admin/getNumBots", cfg.handleGetNumBots)
	mux.HandleFunc("GET /admin/getNumGames", cfg.handleGetNumGames)
	mux.HandleFunc("GET /admin/getCurrentPlayers", cfg.handleGetCurrentPlayers)

	// create the server object
	srv := &http.Server{
		Handler:      mux,
		Addr:         ":" + cfg.port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Printf("server started on http://localhost:%s\n", cfg.port)
	log.Fatal(srv.ListenAndServe()) // this blocks forever, until the server has an unrecoverable error
}
