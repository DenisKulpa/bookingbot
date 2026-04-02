package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLitePath string
	ServerPort string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	sqlitePath := os.Getenv("SQLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "./bookingbot.db"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		SQLitePath: sqlitePath,
		ServerPort: port,
	}, nil
}