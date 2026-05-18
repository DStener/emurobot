package main

import (
	emu "emurobot/shared"
	"os"
	"path/filepath"
	"time"

	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Command struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type Confirmation struct {
	Message string `json:"message"`
}

type DumpInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func hostnameHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Use GET", http.StatusMethodNotAllowed)
		return
	}

	// Config response
	response.Header().Set("Content-Type", "application/json")
	// Send response
	json.NewEncoder(response).Encode(Confirmation{
		Message: EMU_HOSTNAME,
	})
}

func isRecordingHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Use GET", http.StatusMethodNotAllowed)
		return
	}

	// Config response
	response.Header().Set("Content-Type", "application/json")
	// Send response
	json.NewEncoder(response).Encode(Confirmation{
		Message: strconv.FormatBool(isRecording()),
	})
}

func startRecordingHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	log.Info("start record command received")

	dir := filepath.Join(EMU_DUMPS_DIR, time.Now().Format("02.01.2006"), time.Now().Format("15_04_05"))

	setLogPath(dir)

	setRecording(true)
	startRecordingCamera()

	confirmation := Confirmation{
		Message: "Recording started",
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(confirmation)
}

func stopRecordingHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	log.Info("stop record command received")

	dir := getLogPath()

	setRecording(false)
	emu.SaveDumps(dir)

	stopRecordingCamera()

	confirmation := Confirmation{
		Message: "Recording stopped",
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(confirmation)
}

func getDumpsListHandler(response http.ResponseWriter, request *http.Request) {
	// const rootDir = "dumps"

	var results []DumpInfo

	// Check that dir exist
	if _, err := os.Stat(EMU_DUMPS_DIR); os.IsNotExist(err) {
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(results)
		return
	}

	// Going through all folders starting with EMU_DUMPS_DIR
	err := filepath.Walk(EMU_DUMPS_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip root folder
		if path == EMU_DUMPS_DIR {
			return nil
		}

		// Processing only directories
		if info.IsDir() {
			// Checking if there are subdirectories in this directory
			hasSubdirs := false
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				if entry.IsDir() {
					hasSubdirs = true
					break
				}
			}

			// If there are no subdirectories, this is the destination folder
			if !hasSubdirs {
				// Getting the relative path (removing the prefix "dumps/")
				relativePath := path[len(EMU_DUMPS_DIR)+1:]

				// Calculating the folder size
				folderSize, _ := getFolderSize(path)

				results = append(results, DumpInfo{relativePath, folderSize})
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(results)
	return
}

func getFolderSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
