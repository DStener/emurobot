package emuSerial

import (
	"os"
)

// Return environment or default value
func GetEnvOrDefault(key, value string) string {
	out := os.Getenv(key)

	// If not set, return default value
	if out == "" {
		return value
	}

	return out
}
