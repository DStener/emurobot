package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	emurobot "emurobot/shared"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// Get environment variables
var DEVICE_COUNT int = emurobot.GetEnvOrDefault[int]("DEVICE_COUNT", 1)
var DEFAULT_SPEED int = emurobot.GetEnvOrDefault[int]("DEFAULT_SPEED", 9600)

func main() {

	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	// Init device
	for i := 0; i < DEVICE_COUNT; i++ {
		// Create device
		dev := fmt.Sprintf("/dev/ttyUSB%d", i)
		input, _ := emurobot.CreateDevice(dev)

		// Fill random data
		go loopRandomGenerate(input)
	}

	// Run infinite loop
	for {
	}
}

func loopRandomGenerate(dev string) {

	// Wait until devise is not exist
	emurobot.WaitDeviceExist(dev)

	// Config for opening port
	outputConfig := serial.Config{
		Name:        dev,
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

		// Serial-port speed emulation
		emurobot.WaitSend(&outputConfig, count)
	}
}
