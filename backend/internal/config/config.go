package config

import (
	"log"
	"os"
	"time"
)

var (
	Port         string
	Host         string
	CurrentLevel LogLevel
	StartTime    time.Time
	BasePath     = "/api/v1"
)

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func Init() {
	Port = GetEnv("APP_PORT", GetEnv("PORT", "8080"))
	Host = GetEnv("APP_HOST", "0.0.0.0")

	level, err := ParseLogLevel(GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	CurrentLevel = level

	StartTime = time.Now()
}
