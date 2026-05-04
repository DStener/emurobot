package emurobot

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	api "emurobot/pkg/api"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var TIME_DILATION int = api.GetEnvOrDefault[int]("TIME_DILATION", 1)

// Special cmd handler for socat
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
		log.Fatalf("Cmd error: %v \nStderr: %s", err, stderr.String())
	}
}

func CreateDevice(dev string) (string, string) {
	// Name of devices
	input := fmt.Sprintf("%s_GHOST_IN", dev)
	output := dev

	// Configure args
	args := []string{
		fmt.Sprintf("PTY,link=%s,raw,mode=666", input),
		fmt.Sprintf("PTY,link=%s,raw,mode=666", output),
	}

	// Build command
	cmd := exec.Command("socat", args...)
	addKillMeCmd(cmd, []string{input, output}) // See shared/init.go

	// Run handler in background
	go cmdHandler(cmd)

	return input, output
}

func WaitDeviceExist(dev string) error {
	// Timeout
	deadline := time.Now().Add(1 * time.Second)

	// Wait unit device becomes available
	for time.Now().Before(deadline) {
		_, err := os.Stat(dev)
		if err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond) // Pause 100 ms
	}

	return fmt.Errorf("Timeout wait device exist: %s", dev)
}

// Serial-port speed emulation
func WaitSend(port *serial.Config, bytes int) {
	sizeByte := int(port.StopBits) + int(port.Size)
	timePerBit := time.Second / time.Duration(port.Baud)
	delay := timePerBit * time.Duration(sizeByte*bytes) * time.Duration(TIME_DILATION)
	time.Sleep(delay)
}
