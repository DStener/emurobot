package main

import (
	"os"

	emurobot "emurobot/shared"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	// Load config
	path := emurobot.GetEnvOrDefault[string]("EMU_CONFIG_PATH", "/etc/emu_config.yaml")
	config := readConfig(path)

	for _, dev := range config.Devices {
		// Create device
		input, _ := emurobot.CreateDevice(dev.Output)
		dev.GhostInput = input

		// Run main logic
		go runGhostCopy(dev)
	}

	for {

	}
}
