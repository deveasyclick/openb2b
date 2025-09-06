package config

import (
	"fmt"
	"os"

	parseintenv "github.com/deveasyclick/openb2b/internal/utils/parseintEnv"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/joho/godotenv"
)

const (
	defaultPort      = 8080
	defaultDBPort    = 5432
	defaultRedisPort = 6379
	defaultEnv       = "development"
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
	SMTPHost                  string
	SMTPPort                  int
	SMTPUser                  string
	SMTPPassword              string
	SMTPFrom                  string
}

// LoadConfig loads environment variables from .env (if available) and system envs.
func LoadConfig(logger interfaces.Logger) (*Config, error) {
	// Try loading .env files (optional)
	if err := godotenv.Load(".env", ".env.local"); err != nil {
		logger.Warn("No .env file found, falling back to system environment")
	}

	cfg := &Config{
		Env:                       getEnv("ENV", defaultEnv),
		DBHost:                    os.Getenv("DB_HOST"),
		DBName:                    os.Getenv("DB_NAME"),
		DBUser:                    os.Getenv("DB_USER"),
		DBPassword:                os.Getenv("DB_PASSWORD"),
		AppURL:                    os.Getenv("APP_URL"),
		Port:                      parseintenv.ParseIntEnv("PORT", defaultPort, logger),
		DBPort:                    parseintenv.ParseIntEnv("DB_PORT", defaultDBPort, logger),
		RedisPort:                 parseintenv.ParseIntEnv("REDIS_PORT", defaultRedisPort, logger),
		ClerkWebhookSigningSecret: os.Getenv("CLERK_WEBHOOK_SIGNING_SECRET"), // optional
		ClerkSecret:               os.Getenv("CLERK_SECRET_KEY"),
		SMTPHost:                  os.Getenv("SMTP_HOST"),
		SMTPPort:                  parseintenv.ParseIntEnv("SMTP_PORT", 587, logger),
		SMTPUser:                  os.Getenv("SMTP_USER"),
		SMTPPassword:              os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:                  os.Getenv("SMTP_FROM"),
	}

	// Validate required config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Warn about optional but recommended config
	if cfg.ClerkWebhookSigningSecret == "" {
		logger.Warn("Missing optional env variable",
			"key", "CLERK_WEBHOOK_SIGNING_SECRET",
			"reason", "needed for user creation via Clerk webhooks",
		)
	}

	return cfg, nil
}

// Validate ensures all required environment variables are present
func (c *Config) Validate() error {
	required := map[string]string{
		"DB_HOST":          c.DBHost,
		"DB_NAME":          c.DBName,
		"DB_USER":          c.DBUser,
		"DB_PASSWORD":      c.DBPassword,
		"CLERK_SECRET_KEY": c.ClerkSecret,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	if c.SMTPHost == "" || c.SMTPPort == 0 || c.SMTPUser == "" || c.SMTPPassword == "" || c.SMTPFrom == "" {
		return fmt.Errorf("missing required SMTP environment variables")
	}

	return nil
}

// getEnv fetches an environment variable or returns a fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
