package emurobot_api

import (
	"fmt"
	"sync/atomic"
)

const CMD_START_RECORD = "__CMD_START_RECORD__"
const CMD_STOP_RECORD = "__CMD_STOP_RECORD__"
const CMD_PLAY_RECORD = "__CMD_PLAY_RECORD__"

var slots = make(map[string]func(args []string) (string, error))

var FLAG_IS_RECORDING uint32 = 0

func SetRecording(value bool) {
	if value {
		atomic.StoreUint32(&FLAG_IS_RECORDING, 1)
	} else {
		atomic.StoreUint32(&FLAG_IS_RECORDING, 0)
	}
}

func IsRecording() bool {
	return atomic.LoadUint32(&FLAG_IS_RECORDING) == 1
}

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
