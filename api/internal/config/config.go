package config

import (
	"fmt"
	"os"

	parseintenv "github.com/deveasyclick/openb2b/internal/utils/parseintEnv"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                      int
	Env                       string
	DBHost                    string
	DBName                    string
	DBUser                    string
	DBPassword                string
	DBPort                    int
	RedisPort                 int
	AppURL                    string
	ClerkWebhookSigningSecret string
	ClerkSecret               string
}

func LoadConfig(logger interfaces.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{
		Env:                       os.Getenv("ENV"),
		DBHost:                    os.Getenv("DB_HOST"),
		DBName:                    os.Getenv("DB_NAME"),
		DBUser:                    os.Getenv("DB_USER"),
		DBPassword:                os.Getenv("DB_PASSWORD"),
		AppURL:                    os.Getenv("APP_URL"),
		Port:                      parseintenv.ParseIntEnv("PORT", 8080, logger),
		DBPort:                    parseintenv.ParseIntEnv("DB_PORT", 5432, logger),
		RedisPort:                 parseintenv.ParseIntEnv("REDIS_PORT", 6379, logger),
		ClerkWebhookSigningSecret: os.Getenv("CLERK_WEBHOOK_SIGNING_SECRET"),
		ClerkSecret:               os.Getenv("CLERK_SECRET_KEY"),
	}

	if cfg.Env == "" {
		cfg.Env = "development"
	}

	if cfg.DBHost == "" || cfg.DBName == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.ClerkSecret == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	if cfg.ClerkWebhookSigningSecret == "" {
		logger.Warn("missing CLERK_WEBHOOK_SIGNING_SECRET environment variable which is needed for creating users in our system")
	}
	return cfg, nil
}
