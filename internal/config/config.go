package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
	ServerPort    string
	Seed          bool
}

func Load() (*Config, error) {
	_, filename, _, _ := runtime.Caller(0)
	rootEnv := filepath.Join(filepath.Dir(filename), "..", "..", ".env")
	_ = godotenv.Load(".env", "../.env", rootEnv)

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/znimai?sslmode=disable"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		TelegramToken: token,
		DatabaseURL:   databaseURL,
		ServerPort:    port,
		Seed:          os.Getenv("SEED") == "true",
	}, nil
}
