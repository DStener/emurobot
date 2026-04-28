package main

import (
	emuserial "emurobot/internal"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {

	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	// Load config
	path := emuserial.GetEnvOrDefault("EMU_CONFIG_PATH", "/etc/emu_config.yaml")
	config := emuserial.ReadConfig(path)

	emuserial.Init(config)

	for {

	}
}
