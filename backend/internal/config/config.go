package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	Port            string
	Host            string
	CurrentLevel    LogLevel
	StartTime       time.Time
	BasePath        = "/api/v1"
	UsersCollection string
	CardsCollection string
	DecksCollection string
)

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	Port = GetEnv("APP_PORT", GetEnv("PORT", "8080"))
	Host = GetEnv("APP_HOST", "0.0.0.0")
	UsersCollection = GetEnv("USERS_COLLECTION", "users")
	CardsCollection = GetEnv("CARDS_COLLECTION", "cards")
	DecksCollection = GetEnv("DECKS_COLLECTION", "decks")

	level, err := ParseLogLevel(GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	CurrentLevel = level

	StartTime = time.Now()
}
