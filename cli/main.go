package main

import (
	"flag"
	"fmt"

	api "emurobot/pkg/api"

	log "github.com/sirupsen/logrus"
)

func parseArgs() (string, error) {
	var recAction string
	var playPath string

	flag.StringVar(&recAction, "rec", "", "Record: 'start' or 'stop'")
	flag.StringVar(&playPath, "play", "", "Play recorded file")

	flag.Parse()

	if recAction != "" && playPath != "" {
		return "ERROR", fmt.Errorf("You cannot use 'rec' and 'play' at the same time")
	}

	// Run rec
	if recAction != "" {
		switch recAction {
		case "start":
			return api.StartRecord()
		case "stop":
			return api.StopRecord()
		default:
			return "ERROR", fmt.Errorf("for 'rec', only 'start' or 'stop' is allowed")
		}
	}

	// Run play
	if playPath != "" {
		return api.PlayRecord(playPath)
	}

	return "ERROR", fmt.Errorf("the command is not specified")
}

func main() {
	str, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	log.Info(str)
}
