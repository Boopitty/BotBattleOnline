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
			processed := Process(cfg, msg)
			conn.WriteMessage(websocket.TextMessage, processed)
		}
	}
}
