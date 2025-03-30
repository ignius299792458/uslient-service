package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type EnvConfig struct {
	// jwt settings
	JwtScretKey string

	// Database settings
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Application settings
	AppEnv   string
	LogLevel string
	Port     int
}

// Load returns a config struct populated from environment variables
func EnvLoad() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or unable to load")
	}

	port, _ := strconv.Atoi(getEnv("PORT", "9090"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return &EnvConfig{
		// API settings
		JwtScretKey: getEnv("JWT_SCRET_KEY", ""),

		// Database settings
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "myapp"),

		// Application settings
		AppEnv:   getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Port:     port,
	}
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
