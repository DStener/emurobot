package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	api "emurobot/pkg/api"
	emurobot "emurobot/shared"

	log "github.com/sirupsen/logrus"
)

var FLAG_IS_RECORDING uint32 = 0

func SetRecording(value bool) {
	if value {
		atomic.StoreUint32(&FLAG_IS_RECORDING, 1)
	} else {
		atomic.StoreUint32(&FLAG_IS_RECORDING, 0)
	}
}

func IsRecording() bool {
	return atomic.LoadUint32(&FLAG_IS_RECORDING) == 1
}

func StartRecord(args []string) (string, error) {
	log.Info("start record command received")

	if err := callCameraRecorder("/start"); err != nil {
		log.Error("camera-recorder start failed: ", err)
		return "ERROR", err
	}

	SetRecording(true)

	log.Info("recording started")

	return "OK", nil
}

func StopRecord(args []string) (string, error) {
	log.Info("stop record command received")

	SetRecording(false)

	if err := callCameraRecorder("/stop"); err != nil {
		log.Error("camera-recorder stop failed: ", err)
		return "ERROR", err
	}

	emurobot.SaveDumps()

	log.Info("recording stopped")

	return "OK", nil
}

func callCameraRecorder(path string) error {
	baseURL := getEnv("CAMERA_RECORDER_URL", "http://camera-recorder:8080")
	url := strings.TrimRight(baseURL, "/") + path

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("camera-recorder returned status: %s", resp.Status)
	}

	return nil
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}

func main() {
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	go api.InitServer()

	api.Connect(api.CMD_START_RECORD, StartRecord)
	api.Connect(api.CMD_STOP_RECORD, StopRecord)

	path := api.GetEnvOrDefault[string]("EMU_CONFIG_PATH", "/etc/emu_config.yaml")
	config := emurobot.ReadConfig(path)

	for _, dev := range config.Devices {
		input, _ := emurobot.CreateDevice(dev.Output)
		dev.GhostInput = input

		go runGhostCopy(dev, config.BufferSize)
	}

	select {}
}
