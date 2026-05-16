package main

import "sync/atomic"

var FLAG_IS_RECORDING uint32 = 0

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
