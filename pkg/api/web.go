package emurobot_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Command struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type Confirmation struct {
	Message string `json::"message"`
}

const METHOD_PATH = "/api/run"

var EMU_STATIC_PATH string = GetEnvOrDefault[string]("EMU_STATIC_PATH", "static")
var HOST_HOSTNAME string = GetEnvOrDefault[string]("HOST_HOSTNAME", "Robot")
var SIGNAL_ADDRESS string = GetEnvOrDefault[string]("EMU_SIGNAL_ADDRESS", ":7000")

func messageHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var command Command
	err := json.NewDecoder(request.Body).Decode(&command)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	// start message runner
	msg, err := messageRunner(command)

	// Set message status
	if err != nil {
		msg = err.Error()
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.WriteHeader(http.StatusOK)
	}

	// Config response
	response.Header().Set("Content-Type", "application/json")

	confirmation := Confirmation{
		Message: msg,
	}

	// Send response
	json.NewEncoder(response).Encode(confirmation)
}
func SendCommand(command string, args []string) (string, error) {

	// Send to server
	cmd := Command{Cmd: command, Args: args}
	jsonData, err := json.Marshal(cmd)
	if err != nil {
		return "ERROR", err
	}

	resp, err := http.Post(
		SIGNAL_ADDRESS+METHOD_PATH,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "ERROR", fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()

	// Get response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR", err
	}

	var serverConf Confirmation
	err = json.Unmarshal(body, &serverConf)
	if err != nil {
		return "ERROR", err
	}

	if resp.StatusCode != http.StatusOK {
		return "ERROR", fmt.Errorf("RESPONSE ERROR: %s", serverConf.Message)
	}

	return serverConf.Message, nil
}

func hostnameHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Config response
	response.Header().Set("Content-Type", "application/json")
	// Send response
	json.NewEncoder(response).Encode(Confirmation{
		Message: HOST_HOSTNAME,
	})
}

func isRecordingHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Config response
	response.Header().Set("Content-Type", "application/json")
	// Send response
	json.NewEncoder(response).Encode(Confirmation{
		Message: strconv.FormatBool(!IsRecording()),
	})
}

func InitServer() {
	// Set path for command run
	http.HandleFunc(METHOD_PATH, messageHandler)

	log.Printf("Server starting on %s", SIGNAL_ADDRESS)
	log.Fatal(http.ListenAndServe(SIGNAL_ADDRESS, nil))
}

// TODO: Merge with InitServer()
func InitServerWithWeb() {
	// Creating a multiplexer for routing
	mux := http.NewServeMux()

	// Set path for command run
	mux.HandleFunc(METHOD_PATH, messageHandler)

	mux.HandleFunc("/api/hostname", hostnameHandler)
	mux.HandleFunc("/api/rec/status", isRecordingHandler)

	// For static file
	staticHandler := http.StripPrefix("/", http.FileServer(http.Dir(EMU_STATIC_PATH)))
	mux.Handle("/", staticHandler)

	log.Printf("Server starting on %s", SIGNAL_ADDRESS)
	log.Fatal(http.ListenAndServe(SIGNAL_ADDRESS, mux))
}
