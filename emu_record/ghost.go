package main

import (
	"io"
	"time"

	emurobot "emurobot/shared"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func runGhostCopy(dev emurobot.Device, buffer_size int) {

	// Init log structure
	dump := emurobot.InitDevLog(dev)

	// Wait until devise is not exist
	errIn := emurobot.WaitDeviceExist(dev.GhostInput)
	errOut := emurobot.WaitDeviceExist(dev.Output)

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

	// Create buffer
	stream := make([]byte, buffer_size)

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
		start := time.Now()

		// Read from real port
		_, err = input.Read(stream)
		if err != nil {
			// If the buffer does nоt print anything, just wait.
			if err == io.EOF {
				continue
			}
			log.Fatalf("Read error: %s", err)
		}

		log.Info("TEST ", stream)

		// Event logging
		emurobot.AddEvent(
			dump,
			string(stream),
			time.Since(start),
			FLAG_IS_RECORDING,
		)

		// // Write to ghost port
		// _, err = ghost.Write(stream)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

}
