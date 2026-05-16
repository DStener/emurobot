package main

import (
	"io"
	"time"

	emu "emurobot/shared"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func InitSerial(config emu.Config) error {

	// Create serial-port
	for _, dev := range config.Devices {
		input, _ := emu.CreateDevice(dev.Output)
		dev.GhostInput = input

		go startSerialBridge(dev, config.BufferSize)
	}

	return nil
}

func startSerialBridge(dev emu.Device, buffer_size int) {

	// Init log structure
	dump := emu.InitDevLog(dev)

	// Wait until devise is not exist
	errIn := emu.WaitDeviceExist(dev.GhostInput)
	errOut := emu.WaitDeviceExist(dev.Output)

	if errIn != nil || errOut != nil {
		log.Error(errIn, errOut)
	}

	// Create configs
	inputConfig := serial.Config{
		Name:        dev.Input,
		Baud:        dev.Speed,
		ReadTimeout: 1 * time.Second,
		Size:        8,
	}
	ghostConfig := serial.Config{
		Name:        dev.GhostInput,
		Baud:        dev.Speed,
		ReadTimeout: 1 * time.Second,
		Size:        8,
	}

	// Open all ports
	input, err := serial.OpenPort(&inputConfig)
	if err != nil {
		log.Panicln("Not open input port", err.Error())
	}
	defer input.Close()

	ghost, err := serial.OpenPort(&ghostConfig)
	if err != nil {
		log.Panicln("Not open ghost port", err.Error())
	}
	defer ghost.Close()

	// Create buffer
	stream := make([]byte, buffer_size)

	// Main loop
	for {

		start := time.Now()

		var n int

		// Read from real port
		n, err = input.Read(stream)
		if err != nil {
			// If the buffer does nоt print anything, just wait.
			if err == io.EOF {
				continue
			}
			log.Fatalf("Read error: %s", err)
		}

		// log.Info("TEST ", stream)

		// Event logging
		emu.AddEvent(
			dump,
			string(stream[:n]),
			time.Since(start),
			isRecording(),
		)

		// Write to ghost port
		n, err = ghost.Write(stream[:n])
		if err != nil {
			log.Fatal(err)
		}
	}

}
