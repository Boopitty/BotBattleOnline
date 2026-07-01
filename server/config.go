package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// config holds the configuration values for the server.
type config struct {
	port         string
	filepathRoot string
	token        string
}

// create a new config object by reading environment variables
func createConfig() *config {
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

	JWT := os.Getenv("JWT")
	if JWT == "" {
		log.Fatalln("JWT env variable not found.")
	}

	cfg := &config{
		port:         port,
		filepathRoot: filepathRoot,
		token:        JWT,
	}

	return cfg
}
