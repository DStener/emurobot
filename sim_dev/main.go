package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var DEVICE_COUNT int
var DEFAULT_SPEED int

func main() {

	var err1, err2 error

	// Check permission
	if os.Geteuid() != 0 {
		log.Fatal("Root permissions are missing")
	}

	DEVICE_COUNT, err1 = strconv.Atoi(os.Getenv("DEVICE_COUNT"))
	DEFAULT_SPEED, err2 = strconv.Atoi(os.Getenv("DEFAULT_SPEED"))

	if err1 != nil || err2 != nil {
		log.Panic("DEVICE_COUNT or DEFAULT_SPEED incorrect set:", err1, err2)
	}

	// Init device
	for i := 0; i < DEVICE_COUNT; i++ {
		dev := createDevice(i)
		go loopRandomGenerate(dev)
	}

	// Run infinite loop
	for {
	}
}

func createDevice(num int) string {

	// Configure args
	args := []string{
		fmt.Sprintf("PTY,link=/dev/ttySIM%d", num),
		fmt.Sprintf("PTY,link=/dev/ttyUSB%d", num),
	}

	// Build command
	cmd := exec.Command("socat", args...)

	// Run handler in background
	go cmdHandler(cmd)

	return fmt.Sprintf("/dev/ttySIM%d", num)
}

func cmdHandler(cmd *exec.Cmd) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.StdinPipe()

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed: %v \nStderr: %s", err, stderr.String())
	}
}

func loopRandomGenerate(dev string) {

	// Wait until devise is not exist
	waitDevice(dev)

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

	// start-bit + data-bit + stop-bit
	bitsPerByte := int(outputConfig.Size + 2)
	timePerBit := time.Second / time.Duration(DEFAULT_SPEED)

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

		// Delay calc: bits * time per bit
		totalBits := count * bitsPerByte
		delay := timePerBit * time.Duration(totalBits)
		time.Sleep(delay)
	}

}

func waitDevice(dev string) {
	// Timeout
	deadline := time.Now().Add(1 * time.Second)

	// Wait unit device becomes available
	for time.Now().Before(deadline) {
		_, err := os.Stat(dev)
		if err == nil {
			break // Device is exist
		}

		time.Sleep(100 * time.Millisecond) // Pause 100 ms
	}
}
