package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func LoadConfig() *Config {
	// Load environment variables from .env file
	_ = godotenv.Load()

	return &Config{
		Env:           GetEnvOrDefault("APP_ENV", "development"),
		DBHost:        GetEnvOrDefault("DB_HOST", "localhost"),
		DBPort:        GetEnvOrDefault("DB_PORT", "5432"),
		DBUser:        GetEnvOrDefault("DB_USER", "postgres"),
		DBPassword:    GetEnvOrDefault("DB_PASSWORD", ""),
		DBName:        GetEnvOrDefault("DB_NAME", "postgres"),
		RedisHost:     GetEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     GetEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: GetEnvOrDefault("REDIS_PASSWORD", ""),
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
