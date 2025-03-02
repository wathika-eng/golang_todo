// provides secret config keys for our API
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	SERVER_PORT       string
	DB_TYPE           string
	DB_HOST           string
	DB_PORT           string
	DB_USER           string
	DB_PASSWORD       string
	DB_NAME           string
	CONNECTION_STRING string
	SECRET_KEY        string
}

var Envs = initConfigs()

func initConfigs() config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or failed to load. Using default environment variables.")
	}
	return config{
		SERVER_PORT:       getEnv("SERVER_PORT", "8000"),
		DB_TYPE:           getEnv("DB_TYPE", "postgresql"),
		DB_HOST:           getEnv("DB_HOST", "localhost"),
		DB_PORT:           getEnv("DB_PORT", "5432"),
		DB_USER:           getEnv("DB_USER", "postgres"),
		DB_PASSWORD:       getEnv("DB_PASSWORD", "postgres"),
		DB_NAME:           getEnv("DB_NAME", "todoApp"),
		CONNECTION_STRING: getEnv("CONNECTION_STRING", ""),
		SECRET_KEY:        getEnv("SECRET_KEY", "https://acte.ltd/utils/randomkeygen"),
	}
}

// getEnv retrieves an environment variable or returns the fallback string
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
