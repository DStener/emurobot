package emuSerial

import (
	"math"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const CURRENT_VERSION = "1.0.0"

type Config struct {
	Version string   `yaml:"version"`
	Devices []Device `yaml:"devices"`
}

type Device struct {
	Input      string `yaml:"input"`
	GhostInput string
	Output     string `yaml:"output"`
	Speed      int    `yaml:"speed"`
	Size       int    `yaml:"size"`
}

func ReadConfig(path string) Config {
	var out Config

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	// Parse
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		log.Panic(err)
	}

	if !isSupportedVersion(out.Version, CURRENT_VERSION) {
		log.Errorf("Config file version %s not supported (current version %s)", out.Version, CURRENT_VERSION)
	}

	return out
}

func isSupportedVersion(current, required string) bool {
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
