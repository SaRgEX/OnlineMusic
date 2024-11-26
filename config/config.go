package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	HTTPServer
	Database
	LogLevel    string
	LogFilePath string
	ApiURL      string
}

type HTTPServer struct {
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewConfig() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return &Config{
		Database: Database{
			Username: envValue("POSTGRES_USER", "postgres"),
			Host:     envValue("POSTGRES_HOST", "localhost"),
			Port:     envValue("POSTGRES_PORT", "5432"),
			Database: envValue("POSTGRES_DATABASE", "postgres"),
			SSLMode:  envValue("POSTGRES_SSLMODE", "disable"),
			Password: envValue("POSTGRES_PASSWORD", "postgres"),
		},
		HTTPServer: HTTPServer{
			Host: envValue("HTTP_SERVER_ADDRESS", "localhost"),
			Port: envValue("HTTP_SERVER_PORT", ":8080"),
		},
		LogLevel:    envValue("LOG_LEVEL", "info"),
		LogFilePath: envValue("LOG_FILE_PATH", "/log/log.log"),
		ApiURL:      envValue("API_URL", "http://localhost:8000/info"),
	}
}

func envValue(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
