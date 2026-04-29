package emurobot_api

import (
	"math"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Check that current version lower than required
func IsSupportedVersion(current, required string) bool {
	currentParts := strings.Split(current, ".")
	requiredParts := strings.Split(required, ".")

	// Get min len version
	min := int(math.Min(float64(len(currentParts)), float64(len(requiredParts))))

	for i := 0; i < min; i++ {
		currentNum, _ := strconv.Atoi(currentParts[i])
		requiredNum, _ := strconv.Atoi(requiredParts[i])

		if currentNum > requiredNum {
			return false
		}
	}

	return true
}

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
	case string:
		return any(str).(T)
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
