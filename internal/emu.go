package emuSerial

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/tarm/serial"
)

// import

func Init(config Config) {

	// Iteration for create virtual device
	for _, dev := range config.Devices {
		createDevice(&dev)
		go runGhostCopy(dev)
	}
}

// Create virtual device
func createDevice(dev *Device) {

	dev.GhostInput = fmt.Sprintf("%s_GHOST_IN", dev.Output)

	args := []string{
		"-d",
		fmt.Sprintf("pty,raw,echo=0,link=%s,b%d", dev.GhostInput, dev.Speed),
		fmt.Sprintf("pty,raw,echo=0,link=%s,b%d", dev.Output, dev.Speed),
	}

	cmd := exec.Command("socat", args...)

	err := cmd.Start()
	if err != nil {
		log.Fatal("Broken start socat: ", err)
	}

	log.Printf("Device %s is created", dev.Output)
}

func runGhostCopy(dev Device) {

	// Create configs
	inputConfig := serial.Config{
		Name:        dev.Output,
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

		n, err = ghost.Write([]byte("test"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Write %d bytes", n)

		// Read from real port
		n, err = input.Read(stream)
		if err != nil {
			log.Fatal(err)
		}

		// Just pass iteration
		if n == 0 {
			continue
		}

		log.Printf("Reade %d byte: %q", n, stream[:n])

		// // Write to ghost port
		// n, err = ghost.Write(stream[:n])
		// if err != nil {
		// 	log.Fatal(err)
		// }

	}
}
