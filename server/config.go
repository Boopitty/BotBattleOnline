package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Boopitty/BotBattleOnline/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// config holds the configuration values for the server.
type config struct {
	port         string
	filepathRoot string
	secret       string
	db           *database.Queries
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

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatalln("JWT env variable not found.")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL env variable not found")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	cfg := &config{
		port:         port,
		filepathRoot: filepathRoot,
		secret:       secret,
		db:           database.New(dbConn),
	}
	return cfg
}
