package utils

import (
	"log/slog"
	"os"
	"strconv"
)

func ParseIntEnv(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		slog.Warn("Invalid integer for "+key+", using default", "err", err, "default", defaultVal)
		return defaultVal
	}
	return val
}
