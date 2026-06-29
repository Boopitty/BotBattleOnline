package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Server")
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load godotenv %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("PORT env variable not found.")
	}

	filepathRoot := os.Getenv("FILEPATH_ROOT")
	if filepathRoot == "" {
		log.Fatalln("FILEPATH_ROOT env variable not found.")
	}

	// Mutex for the server
	mux := http.NewServeMux()
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	// Handle the front page
	mux.Handle("/app/", appHandler)

	// create the server object
	srv := http.Server{
		Handler:      mux,
		Addr:         ":" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// this blocks forever, until the server
	// has an unrecoverable error
	fmt.Printf("server started on %v\n", port)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
