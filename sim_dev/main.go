package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	emurobot "emurobot/shared"

	api "emurobot/pkg/api"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// Get environment variables
var DEVICE_COUNT int = api.GetEnvOrDefault[int]("DEVICE_COUNT", 1)
var DEFAULT_SPEED int = api.GetEnvOrDefault[int]("DEFAULT_SPEED", 9600)

func main() {
	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	// Init device
	for i := 0; i < DEVICE_COUNT; i++ {

		dev := fmt.Sprintf("/dev/ttyUSB%d", i)
		duplicate := fmt.Sprintf("/dev/ttyUSB%d_DEBUG", i)

		// Create devices
		input, _ := emurobot.CreateDevice(dev)
		input_dupl, _ := emurobot.CreateDevice(duplicate)

		// Fill random data
		go loopRandomGenerate(input, input_dupl)
	}

	// Run infinite loop
	for {
	}
}

func loopRandomGenerate(dev string, dupl string) {

	// Wait until devise is not exist
	emurobot.WaitDeviceExist(dev)
	emurobot.WaitDeviceExist(dupl)

	// Config for opening port
	outputConfig := serial.Config{
		Name:        dev,
		Baud:        DEFAULT_SPEED,
		ReadTimeout: 1 * time.Second,
		Size:        8,
	}
	duplConfig := serial.Config{
		Name:        dupl,
		Baud:        DEFAULT_SPEED,
		ReadTimeout: 1 * time.Second,
		Size:        8,
	}

	// Open port
	output, err := serial.OpenPort(&outputConfig)
	if err != nil {
		log.Panic("Not open input port", err.Error())
	}
	defer output.Close()

	output_dupl, err := serial.OpenPort(&duplConfig)
	if err != nil {
		log.Panic("Not open input port", err.Error())
	}
	defer output.Close()

	// Main loop
	for {

		// Printable buffer
		var buffer []byte

		// Generate count of char
		count := rand.Intn(12)

		// Fill buffer randoms byte
		for i := 0; i < count; i++ {
			buffer = append(buffer, byte(rand.Intn(90)+32))
		}

		// Write buffer
		_, err = output.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}
		_, err = output_dupl.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// Serial-port speed emulation
		emurobot.WaitSend(&outputConfig, count)
	}
}
