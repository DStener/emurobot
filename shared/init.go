package emurobot

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

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

var cmdList []*exec.Cmd

// Add cmd to the sheet to be destroyed
func addKillMeCmd(cmd *exec.Cmd) {
	listRWMutex.Lock()
	cmdList = append(cmdList, cmd)
	listRWMutex.Unlock()
}

var once sync.Once

func killAll() {

	listRWMutex.RLock()
	for _, cmd := range cmdList {
		log.Debugf("Kill cmd: %d", cmd.Process.Pid)
		cmd.Process.Signal(syscall.SIGINT)
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
