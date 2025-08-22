// parseintenv parse string to integer from environment variable
//
// It takes a key and a default value. If the environment variable is not set, it returns the default value.

package parseintenv

import (
	"log/slog"
	"os"
	"strconv"
)

// ParseIntEnv parse string to integer from environment variable
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
