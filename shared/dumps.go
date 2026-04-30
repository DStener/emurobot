package emurobot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	api "emurobot/pkg/api"

	log "github.com/sirupsen/logrus"
)

type LogDump struct {
	Device string     `json:"dev"`
	Speed  int        `json:"speed"`
	Size   int        `json:"size"`
	Events []LogEvent `json:"events"`
}

type LogEvent struct {
	Time  uint32 `json:"time"`
	Bytes []byte `json:"bytes"`
}

var dumps []*LogDump

var DUMPS_DIR = api.GetEnvOrDefault[string]("EMU_DUMPS_DIR", "/var/log/emudump")

func InitDevLog(dev string, speed, size int) *LogDump {

	out := LogDump{
		Device: dev,
		Speed:  speed,
		Size:   size,
		Events: make([]LogEvent, 0),
	}

	dumps = append(dumps, &out)

	return &out
}

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
