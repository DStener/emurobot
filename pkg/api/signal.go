package emurobot_api

import "fmt"

const CMD_START_RECORD = "__CMD_START_RECORD__"
const CMD_STOP_RECORD = "__CMD_STOP_RECORD__"
const CMD_PLAY_RECORD = "__CMD_PLAY_RECORD__"

var slots = make(map[string]func(args []string) (string, error))

// Add slot for this signal
func Connect(signal string, slot func(args []string) (string, error)) {
	slots[signal] = slot
}

func messageRunner(command Command) (string, error) {

	if slot, exists := slots[command.Cmd]; exists {
		return slot(command.Args)
	}

	return "ERROR", fmt.Errorf("No bind slots for command: %s", command.Cmd)
}
