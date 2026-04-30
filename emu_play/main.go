package main

import (
	"fmt"
	"os"

	api "emurobot/pkg/api"
	emurobot "emurobot/shared"

	log "github.com/sirupsen/logrus"
)

var Config emurobot.Config

func PlayRecord(args []string) (string, error) {
	if len(args) < 1 {
		return "ERROR", fmt.Errorf("Incorrect args:", args)
	}
	return InitPlayer(args[0])
}

func main() {
	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	go api.InitServer()

	// Connect to signals
	api.Connect(api.CMD_PLAY_RECORD, PlayRecord)

	// Load config
	path := api.GetEnvOrDefault[string]("EMU_CONFIG_PATH", "/etc/emu_config.yaml")
	Config = emurobot.ReadConfig(path)

	for _, dev := range Config.Devices {
		// Create device
		input, _ := emurobot.CreateDevice(dev.Output)
		dev.GhostInput = input
	}

	for {

	}
}
