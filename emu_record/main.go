package main

import (
	"os"

	api "emurobot/pkg/api"
	emurobot "emurobot/shared"

	log "github.com/sirupsen/logrus"
)

var FLAG_IS_RECORDING bool = false

func StartRecord(args []string) (string, error) {

	FLAG_IS_RECORDING = true

	return "OK", nil
}

func StopRecord(args []string) (string, error) {

	FLAG_IS_RECORDING = false

	for _, dump := range logs {
		log.Info(dump)
	}

	return "OK", nil
}

func main() {

	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	go api.InitServer()

	// Connect to signals
	api.Connect(api.CMD_START_RECORD, StartRecord)
	api.Connect(api.CMD_STOP_RECORD, StopRecord)

	// Load config
	path := api.GetEnvOrDefault[string]("EMU_CONFIG_PATH", "/etc/emu_config.yaml")
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
