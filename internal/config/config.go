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
	SQLitePath    string
	ServerPort    string
}

func Load() (*Config, error) {
	_, filename, _, _ := runtime.Caller(0)
	rootEnv := filepath.Join(filepath.Dir(filename), "..", "..", ".env")
	_ = godotenv.Load(".env", "../.env", rootEnv)

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	sqlitePath := os.Getenv("SQLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "./bookingbot.db"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		TelegramToken: token,
		SQLitePath:    sqlitePath,
		ServerPort:    port,
	}, nil
}
