package config

import (
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	ServerPort string
	DBPath     string
	UploadPath string
	LogPath    string
}

func LoadConfig() *Config {
	uploadPath := filepath.Join("data", "uploads")
	logPath := filepath.Join("logs")

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
	}

	if err := os.MkdirAll(logPath, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBPath:     filepath.Join("data", "obsync.db"),
		UploadPath: uploadPath,
		LogPath:    logPath,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
