package main

import (
	"sync"
	"sync/atomic"
)

var FLAG_IS_RECORDING uint32 = 0

var CURRENT_LOG_PATH string = ""

var logPathMutex sync.RWMutex

func setRecording(value bool) {
	if value {
		atomic.StoreUint32(&FLAG_IS_RECORDING, 1)
	} else {
		atomic.StoreUint32(&FLAG_IS_RECORDING, 0)
	}
}

func isRecording() bool {
	return atomic.LoadUint32(&FLAG_IS_RECORDING) == 1
}

func setLogPath(value string) {
	logPathMutex.Lock()
	CURRENT_LOG_PATH = value
	logPathMutex.Unlock()
}

func getLogPath() string {
	logPathMutex.RLock()
	out := CURRENT_LOG_PATH
	logPathMutex.RUnlock()

	return out
}
