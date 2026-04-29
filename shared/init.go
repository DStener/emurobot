package emurobot

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

/* When "docker compose down" is executed or SIGINT and SIGTERM
 * signals are sent, not all instances of cmd socat have time
 * to stop and the viral serial ports stop hanging in the system.
 * Therefore, before finally executing os.Exit(0), it is necessary
 * to manually send SIGINT to each cmd
 */

// Mutex for work with cmdList
var listRWMutex sync.RWMutex

type Record struct {
	cmd  *exec.Cmd
	devs []string
}

var cmdList []Record

var once sync.Once

// Add cmd to the sheet to be destroyed
func addKillMeCmd(cmd *exec.Cmd, files []string) {
	listRWMutex.Lock()
	cmdList = append(cmdList, Record{cmd, files})
	listRWMutex.Unlock()

}

func waitNonExist(dev string) {
	deadline := time.Now().Add(1 * time.Second)

	// Wait unit device becomes available
	for time.Now().Before(deadline) {
		_, err := os.Stat(dev)
		if err != nil {
			break // Device is non-exist
		}

		time.Sleep(100 * time.Millisecond) // Pause 100 ms
	}
}

func killAll() {

	listRWMutex.RLock()

	for _, rec := range cmdList {
		log.Debugf("Kill cmd: %d", rec.cmd.Process.Pid)
		rec.cmd.Process.Signal(syscall.SIGINT)

		// Wait while device is deleted
		for _, dev := range rec.devs {
			waitNonExist(dev)
		}
	}
	listRWMutex.RUnlock()

	os.Exit(0)
}

func waitKill() {
	// Create chan for signal
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Infinite loop
	for {
		<-sigChan
		killAll()
	}
}

func init() {

	// Set correct exit handler
	log.RegisterExitHandler(killAll)

	go waitKill()
}
