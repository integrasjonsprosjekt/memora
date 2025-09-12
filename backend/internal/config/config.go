package config

import (
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

var (
	Port           string
	Host           string
	CurrentLevel   LogLevel
	StartTime      time.Time
	BasePath       = "/api/v1"
	FirebaseClient *firestore.Client
)

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	Port = GetEnv("APP_PORT", GetEnv("PORT", "8080"))
	Host = GetEnv("APP_HOST", "0.0.0.0")

	level, err := ParseLogLevel(GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	CurrentLevel = level

	StartTime = time.Now()
}
