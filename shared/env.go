package emurobot

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Return environment or default value
func GetEnvOrDefault[T any](key string, value T) T {
	var out T

	// Get value of environment
	str := os.Getenv(key)

	// If not set, return default value
	if str == "" {
		return value
	}

	// Convert to support type
	switch any(out).(type) {
	case int:
		i, err := strconv.Atoi(str)
		if err != nil {
			log.Errorf("Error convert: %s", err)
		}
		return any(i).(T)

	case float32:
		f, err := strconv.ParseFloat(str, 32)
		if err != nil {
			log.Errorf("Error convert: %s", err)
		}
		return any(f).(T)

	case float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Errorf("Error convert: %s", err)
		}
		return any(f).(T)

	case bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			log.Errorf("Error convert: %s", err)
		}
		return any(b).(T)
	}

	// Error if type is unsupported
	log.Errorf("Unsupported type: %T", value)
	return any(str).(T)
}
