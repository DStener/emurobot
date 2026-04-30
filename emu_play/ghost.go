package main

import (
	emurobot "emurobot/shared"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func InitPlayer(path string) (string, error) {

	err := emurobot.LoadDumps(path)
	if err != nil {
		return "ERROR", err
	}

	for i := 0; i < emurobot.GetDeviceCount(); i++ {
		dump := emurobot.GetDeviceLog(i)
		go runGhostPlayer(*dump)
	}

	return "START", nil
}

func runGhostPlayer(dump emurobot.LogDump) {

	// Wait until devise is not exist
	errIn := emurobot.WaitDeviceExist(dump.Input)
	errOut := emurobot.WaitDeviceExist(dump.Output)

	if errIn != nil || errOut != nil {
		log.Error(errIn, errOut)
	}

	// Create configs
	inputConfig := serial.Config{
		Name:        dump.Input,
		Baud:        dump.Speed,
		ReadTimeout: 1 * time.Second,
		Size:        byte(dump.Size),
	}

	input, err := serial.OpenPort(&inputConfig)
	if err != nil {
		log.Panicln("Not open input port", err.Error())
	}
	defer input.Close()

	// Loop for simulate real device
	for _, event := range dump.Events {
		// Simulate input time
		time.Sleep(time.Duration(event.Time))

		// Write to ghost port
		_, err := input.Write([]byte(event.Data))
		if err != nil {
			log.Fatal(err)
		}
	}
}
