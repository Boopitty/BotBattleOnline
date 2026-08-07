package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := createConfig()

	mux := http.NewServeMux()                                 // Mutex for the server
	appHandler := http.FileServer(http.Dir(cfg.filepathRoot)) // Serve the static files from the specified root directory

	// Handle the front page
	mux.Handle("/", appHandler)

	// Handle the API endpoints
	mux.HandleFunc("/api/command", cfg.handleCommand)
	mux.HandleFunc("POST /api/createUser", cfg.handleCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)
	mux.HandleFunc("DELETE /api/deleteUser", cfg.handleDeleteUser)
	mux.HandleFunc("DELETE /api/resetUsers", cfg.handleReset)
	mux.HandleFunc("/ws", handleWS(cfg))

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
