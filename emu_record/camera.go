package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	emu "emurobot/shared"
)

const (
	inputCamera   = "/dev/video0"
	virtualCamera = "/dev/video10"

	videoSize   = "640x480"
	fps         = "5"
	inputFormat = "yuyv422"
)

var (
	mu sync.Mutex

	bridgeCmd *exec.Cmd

	recordCmd  *exec.Cmd
	recordFile string
	recordDone chan error
)


func callCameraRecorder(path string) error {
	baseURL := emu.GetEnv("CAMERA_RECORDER_URL", "http://camera-recorder:8080")

	url := strings.TrimRight(baseURL, "/") + path

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("camera-recorder returned status: %s", resp.Status)
	}

	return nil
}

func startBridge() error {
	if bridgeCmd != nil {
		log.Println("camera bridge already running")
		return nil
	}

	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-f", "v4l2",
		"-framerate", fps,
		"-video_size", videoSize,
		"-i", inputCamera,

		"-vf", "format=yuyv422",
		"-vcodec", "rawvideo",
		"-pix_fmt", "yuyv422",
		"-f", "v4l2",
		virtualCamera,
	}

	cmd := exec.Command("ffmpeg", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	bridgeCmd = cmd

	log.Println("camera bridge started:", inputCamera, "->", virtualCamera)

	go func() {
		err := cmd.Wait()

		if bridgeCmd == cmd {
			bridgeCmd = nil
		}

		log.Println("camera bridge stopped:", err)
	}()

	return nil
}

func stopBridge() error {
	if bridgeCmd == nil {
		log.Println("camera bridge is not running")
		return nil
	}

	cmd := bridgeCmd
	bridgeCmd = nil

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		return syscall.Kill(-pgid, syscall.SIGINT)
	}

	return cmd.Process.Kill()
}

func startRecordingCamera() {

	mu.Lock()
	defer mu.Unlock()

	if recordCmd != nil {
		// fmt.Fprintln(w, "recording already started")
		return
	}

	dir := getLogPath()

	if err := os.MkdirAll(dir, 0755); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recordFile = fmt.Sprintf(
		"%s/record_%s.mp4",
		dir,
		time.Now().Format("2006-01-02_15-04-05"),
	)

	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-y",

		"-f", "v4l2",
		"-input_format", inputFormat,
		"-framerate", fps,
		"-video_size", videoSize,
		"-i", virtualCamera,

		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-pix_fmt", "yuv420p",

		recordFile,
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recordCmd = cmd
	recordDone = make(chan error, 1)

	log.Println("recording started:", recordFile)

	go waitRecording(cmd)

	// fmt.Fprintln(w, "recording started")
}

func stopRecordingCamera() {

	mu.Lock()

	if recordCmd == nil {
		mu.Unlock()
		// fmt.Fprintln(w, "recording is not running")
		return
	}

	cmd := recordCmd
	file := recordFile
	done := recordDone

	mu.Unlock()

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGINT)
	} else {
		_ = cmd.Process.Signal(os.Interrupt)
	}

	select {
	case <-done:
		log.Println("recording stopped:", file)
		// fmt.Fprintln(w, "recording stopped")

	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		log.Println("recording killed by timeout:", file)
		// fmt.Fprintln(w, "recording killed by timeout")
	}
}

func waitRecording(cmd *exec.Cmd) {
	err := cmd.Wait()

	mu.Lock()
	defer mu.Unlock()

	if recordCmd == cmd {
		recordCmd = nil
	}

	if err != nil {
		log.Println("recording ffmpeg stopped with error:", err)
	} else {
		log.Println("recording ffmpeg finished normally")
	}

	if recordDone != nil {
		recordDone <- err
		close(recordDone)
	}
}
