package emurobot_api

// Call start record serial-ports (ONLY FOR emu_record)
func StartRecord() (string, error) {
	return SendCommand(CMD_START_RECORD, []string{})
}

// Call stop record serial-ports (ONLY FOR emu_record)
func StopRecord() (string, error) {
	return SendCommand(CMD_STOP_RECORD, []string{})
}

// Call play record (ONLY FOR emu_player)
func PlayRecord(path string) (string, error) {
	return SendCommand(CMD_PLAY_RECORD, []string{path})
}
