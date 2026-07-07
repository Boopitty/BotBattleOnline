package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("Starting Server")

	cfg := createConfig()

	// Mutex for the server
	mux := http.NewServeMux()
	appHandler := http.FileServer(http.Dir(cfg.filepathRoot))

	// Handle the front page
	mux.Handle("/", appHandler)

	mux.HandleFunc("/api/command", cfg.handleCommand)
	mux.HandleFunc("/ws", handleWS(cfg))

	// create the server object
	srv := &http.Server{
		Handler:      mux,
		Addr:         ":" + cfg.port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// this blocks forever, until the server
	// has an unrecoverable error
	fmt.Printf("server started on http://localhost:%s\n", cfg.port)
	log.Fatal(srv.ListenAndServe())
}
