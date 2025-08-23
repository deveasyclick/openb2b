package config

import (
	"fmt"
	"os"

	parseintenv "github.com/deveasyclick/openb2b/internal/utils/parseintEnv"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
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

func LoadConfig(logger interfaces.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{
		Env:        os.Getenv("ENV"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		AppURL:     os.Getenv("APP_URL"),
		Port:                      parseintenv.ParseIntEnv("PORT", 8080, logger),
		DBPort:                    parseintenv.ParseIntEnv("DB_PORT", 5432, logger),
		RedisPort:                 parseintenv.ParseIntEnv("REDIS_PORT", 6379, logger),
	}

	if cfg.Env == "" {
		cfg.Env = "development"
	}

	if cfg.DBHost == "" || cfg.DBName == "" || cfg.DBUser == "" || cfg.DBPassword == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	return cfg, nil
}
