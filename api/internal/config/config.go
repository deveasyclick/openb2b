package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/deveasyclick/openb2b/internal/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       int
	Env        string
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     int
	RedisPort  int
	AppURL     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{
		Port:       utils.ParseIntEnv("PORT", 8080),
		Env:        os.Getenv("ENV"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBPort:     utils.ParseIntEnv("DB_PORT", 5432),
		RedisPort:  utils.ParseIntEnv("REDIS_PORT", 6379),
		AppURL:     os.Getenv("APP_URL"),
	}

	if cfg.Env == "" {
		cfg.Env = "development"
	}

	if cfg.DBHost == "" || cfg.DBName == "" || cfg.DBUser == "" || cfg.DBPassword == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	return cfg, nil
}
