package emurobot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	api "emurobot/pkg/api"

	rb "github.com/floscodes/golang-ringbuffer"
	log "github.com/sirupsen/logrus"
)

type LogDump struct {
	Input  string     `json:"input"`
	Output string     `json:"output"`
	Speed  int        `json:"speed"`
	Size   int        `json:"size"`
	Events []LogEvent `json:"events"`
}

type LogEvent struct {
	Time time.Duration `json:"time"`
	Data string        `json:"data"`
}

var dumps []*LogDump
var buffers map[string]*rb.RingBuffer

var DUMPS_DIR = api.GetEnvOrDefault[string]("EMU_DUMPS_DIR", "/var/log/emudump")
var BUFFER_SIZE = api.GetEnvOrDefault[uint]("EMU_BUFFER_SIZE", 5)

// Init structure for logs
func InitDevLog(dev Device) *LogDump {
	// Fill LogDump structure
	logDump := LogDump{
		Input:  dev.GhostInput,
		Output: dev.Output,
		Speed:  dev.Speed,
		Size:   dev.Size,
		Events: make([]LogEvent, 0),
	}

	// Create ring buffer
	buff := rb.New(BUFFER_SIZE)

	dumps = append(dumps, &logDump)
	buffers[dev.GhostInput] = &buff

	return &logDump
}

// Save dump to file
func SaveDumps() {
	// Get dum dir
	file := fmt.Sprintf("%s_dumps_SPs.json", time.Now().Format("15_04_05"))
	dir := filepath.Join(DUMPS_DIR, time.Now().Format("02.01.2006"))

	// Create dump dir
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error("Dont create dir:", err)
	}

	// Marshal dumps
	data, err := json.Marshal(dumps)
	if err != nil {
		log.Errorf("JSON Marshal error: %s", err)
	}

	// Open file
	err = os.WriteFile(filepath.Join(dir, file), []byte(data), 0644)
	if err != nil {
		log.Error("Bad  attempt write to file:", err)
	}

}

// Load dumps form current file
func LoadDumps(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse
	err = json.Unmarshal(data, &dumps)
	if err != nil {
		return err
	}

	return nil
}

// Return device count
func GetDeviceCount() int {
	return len(dumps)
}

// Return point to logs device row
func GetDeviceLog(index int) *LogDump {
	if index >= len(dumps) {
		log.Errorf("No device with %d index (max %d)", index, len(dumps))
	}
	return dumps[index]
}

func AddEvent(dump *LogDump, data string, time time.Duration, flag_to_log bool) error {
	// Fill structure
	event := LogEvent{
		Time: time,
		Data: data,
	}

	// Check that device is exit
	buf, ok := buffers[dump.Input]
	if !ok {
		return fmt.Errorf("Log Structure for device %s not init", dump.Input)
	}

	// Add event to buffer
	buf.Push(event)

	// If there really is real logging going on
	if flag_to_log {
		// If it first event, fill log from buffer
		if len(dump.Events) == 0 {
			for {
				if buf.IsEmpty() {
					break
				}
				buf_dump, _ := buf.Pop()
				dump.Events = append(dump.Events, any(buf_dump).(LogEvent))
			}
		}
		dump.Events = append(dump.Events, event)
	}

	return nil
}
