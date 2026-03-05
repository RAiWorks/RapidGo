package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Load reads the .env file and sets environment variables.
// If no .env file is found, it logs a message and continues
// (system environment variables are still available).
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}
}
