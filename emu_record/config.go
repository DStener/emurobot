package main

import (
	"os"

	api "emurobot/pkg/api"

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

func readConfig(path string) Config {
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

	if !api.IsSupportedVersion(out.Version, CURRENT_VERSION) {
		log.Errorf("Config file version %s not supported (current version %s)", out.Version, CURRENT_VERSION)
	}

	return out
}
