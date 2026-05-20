package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Confirmation struct {
	Message string `json:"message"`
}

func startPlayerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	videoPath := request.URL.Query().Get("path")
	if videoPath == "" {
		http.Error(response, "Missing video path", http.StatusBadRequest)
		return
	}

	log.Info("start player command received")

	status, err := InitCameraPlayer(videoPath)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	confirmation := Confirmation{
		Message: status,
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(confirmation)
}

func stopPlayerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	log.Info("stop player command received")

	if err := StopCameraPlayer(); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	confirmation := Confirmation{
		Message: "Player stopped",
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(confirmation)
}
