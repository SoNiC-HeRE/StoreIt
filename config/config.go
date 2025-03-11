package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
// If the .env file is not found, it logs a message and falls back to system environment variables.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}
}

// GetMongoURI returns the MongoDB connection URI from the environment variables.
func GetMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set in environment variables")
	}
	return uri
}
