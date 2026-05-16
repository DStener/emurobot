package main

import (
	"net/http"
	"os"

	emu "emurobot/shared"

	log "github.com/sirupsen/logrus"
)

var EMU_CONFIG_PATH = emu.GetEnv("EMU_CONFIG_PATH", "/etc/emurobot/config.yaml")
var EMU_STATIC_PATH = emu.GetEnv("EMU_STATIC_PATH", "static")
var EMU_HOSTNAME = emu.GetEnv("EMU_HOSTNAME", "Robot")
var EMU_ADDRESS = emu.GetEnv("EMU_ADDRESS", ":7000")
var EMU_DUMPS_DIR = emu.GetEnv("EMU_DUMPS_DIR", "/var/log/emudump")

func main() {
	// Check root permissions
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	config := emu.ReadConfig(EMU_CONFIG_PATH)

	// Init serial-port infrastructure
	if err := InitSerial(config); err != nil {
		log.Fatal(err)
	}

	// Init camera infrastructure
	if err := InitCamera(config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/hostname", hostnameHandler)
	http.HandleFunc("/api/dumps", getDumpsListHandler)

	http.HandleFunc("/api/rec/start", startRecordingHandler)
	http.HandleFunc("/api/rec/stop", stopRecordingHandler)
	http.HandleFunc("/api/rec/status", isRecordingHandler)

	// For static file
	staticHandler := http.StripPrefix("/", http.FileServer(http.Dir(EMU_STATIC_PATH)))
	http.Handle("/", staticHandler)

	log.Printf("Server starting on %s", EMU_ADDRESS)
	log.Fatal(http.ListenAndServe(EMU_ADDRESS, nil))

	// select {}
}
