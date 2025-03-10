package config

import (
    "os"
    "github.com/joho/godotenv"
    "log"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }
}

func GetMongoURI() string {
    return os.Getenv("MONGO_URI")
}
