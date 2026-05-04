package emurobot

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const CURRENT_VERSION = "1.0.0"

type GlobalConfig struct {
	Config Config `yaml:"emurobot"`
}

type Config struct {
	Version    string   `yaml:"version"`
	BufferSize int      `yaml:"buffer_size"`
	Devices    []Device `yaml:"serial"`
}

type Device struct {
	Input      string `yaml:"input"`
	GhostInput string
	Output     string `yaml:"output"`
	Speed      int    `yaml:"speed"`
	Size       int    `yaml:"size"`
}

func ReadConfig(path string) Config {
	var cfg GlobalConfig

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	// Parse
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Panic(err)
	}

	return cfg.Config
}
