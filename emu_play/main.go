package main

import (
	"fmt"
	// api "emurobot/pkg/api"
	emu "emurobot/shared"
	emurobot "emurobot/shared"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var EMU_CONFIG_PATH = emu.GetEnv("EMU_CONFIG_PATH", "/etc/emurobot/config.yaml")
var EMU_STATIC_PATH = emu.GetEnv("EMU_STATIC_PATH", "static")
var EMU_HOSTNAME = emu.GetEnv("EMU_HOSTNAME", "Robot")
var EMU_ADDRESS = emu.GetEnv("EMU_ADDRESS", ":7000")
var EMU_DUMPS_DIR = emu.GetEnv("EMU_DUMPS_DIR", "/var/log/emudump")

var Config emurobot.Config

func PlayRecord(args []string) (string, error) {
	if len(args) < 1 {
		return "ERROR", fmt.Errorf("Incorrect args:", args)
	}
	return InitPlayer(args[0])
}

func main() {

	// Init camera infrastructure
	if err := emurobot.InitCamera(Config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/player/start", startPlayerHandler)
	http.HandleFunc("/api/player/stop", stopPlayerHandler)

	log.Printf("Server starting on %s", EMU_ADDRESS)
	log.Fatal(http.ListenAndServe(EMU_ADDRESS, nil))
}
