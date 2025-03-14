// provides secret config keys for our API
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	ServerPort       string
	DbType           string
	DbHost           string
	DbPort           string
	DbUser           string
	DbPassword       string
	DbName           string
	ConnectionString string
	SecretKey        string
	RefreshKey       string
	ResendApiKey     string
	UptraceDsn       string
	GinMode          string
	RedisUrl         string
}

var Envs = initConfigs()

func initConfigs() config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or failed to load. Using default environment variables.")
	}
	return config{
		ServerPort:       getEnv("SERVER_PORT", "8000"),
		DbType:           getEnv("DB_TYPE", "postgresql"),
		DbHost:           getEnv("DB_HOST", "localhost"),
		DbPort:           getEnv("DB_PORT", "5432"),
		DbUser:           getEnv("DB_USER", "postgres"),
		DbPassword:       getEnv("DB_PASSWORD", "postgres"),
		DbName:           getEnv("DB_NAME", "todoApp"),
		ConnectionString: getEnv("CONNECTION_STRING", ""),
		SecretKey:        getEnv("SECRET_KEY", "https://acte.ltd/utils/randomkeygen"),
		RefreshKey:       getEnv("REFRESH_KEY", "https://randomkeygen.com/"),
		ResendApiKey:     getEnv("RESEND_API_KEY", ""),
		UptraceDsn:       getEnv("UPTRACE_DSN", ""),
		GinMode:          getEnv("GIN_MODE", "release"),
		RedisUrl:         getEnv("REDIS_URL", "localhost:6379"),
	}
}

// getEnv retrieves an environment variable or returns the fallback string
func getEnv(key, fallback string) string {
	envKey, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%v key not found, using default value", envKey)
		return fallback
	}
	return envKey
}
