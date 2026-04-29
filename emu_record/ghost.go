package main

import (
	emurobot "emurobot/shared"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func runGhostCopy(dev Device) {

	// Wait until devise is not exist
	emurobot.WaitDeviceExist(dev.GhostInput)
	emurobot.WaitDeviceExist(dev.Output)

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

	// Create buffer
	stream := make([]byte, dev.Size*2)

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

	// Main loop
	for {
		var n int

		// Read from real port
		n, err = input.Read(stream)

		if err != nil && err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Read error: %s", err)
		}

		// Just pass iteration
		if n == 0 {
			continue
		}

		// Write to ghost port
		n, err = ghost.Write(stream[:n])
		if err != nil {
			log.Fatal(err)
		}
	}

}
