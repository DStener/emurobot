package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

const (
	cameraDevice = "/dev/video10"
	recordDir    = "/recordings"

	videoSize   = "640x480"
	fps         = "5"
	inputFormat = "yuyv422"
)

var (
	mu         sync.Mutex
	recordCmd *exec.Cmd
	recordFile string
	done       chan error
)

func main() {
	http.HandleFunc("/start", startRecording)
	http.HandleFunc("/stop", stopRecording)
	http.HandleFunc("/state", state)

	log.Println("camera-recorder started on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func startRecording(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if recordCmd != nil {
		fmt.Fprintln(w, "recording already started")
		return
	}

	if err := os.MkdirAll(recordDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recordFile = fmt.Sprintf(
		"%s/record_%s.mp4",
		recordDir,
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
		"-i", cameraDevice,

		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-pix_fmt", "yuv420p",

		recordFile,
	}

	cmd := exec.Command("ffmpeg", args...)

	// Нужна отдельная process group,
	// чтобы корректно остановить ffmpeg через SIGINT.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recordCmd = cmd
	done = make(chan error, 1)

	log.Println("recording started:", recordFile)

	go waitFFmpeg(cmd)

	fmt.Fprintln(w, "recording started")
}

func stopRecording(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()

	if recordCmd == nil {
		mu.Unlock()
		fmt.Fprintln(w, "recording is not running")
		return
	}

	cmd := recordCmd
	file := recordFile
	waitDone := done

	mu.Unlock()

	// SIGINT нужен, чтобы ffmpeg нормально закрыл mp4-файл.
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGINT)
	} else {
		_ = cmd.Process.Signal(os.Interrupt)
	}

	select {
	case <-waitDone:
		log.Println("recording stopped:", file)
		fmt.Fprintln(w, "recording stopped")

	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		log.Println("recording killed by timeout:", file)
		fmt.Fprintln(w, "recording killed by timeout")
	}
}

func state(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if recordCmd == nil {
		fmt.Fprintln(w, "stopped")
		return
	}

	fmt.Fprintln(w, "recording")
}

func waitFFmpeg(cmd *exec.Cmd) {
	err := cmd.Wait()

	mu.Lock()
	defer mu.Unlock()

	if recordCmd == cmd {
		recordCmd = nil
	}

	if err != nil {
		log.Println("ffmpeg stopped with error:", err)
	} else {
		log.Println("ffmpeg finished normally")
	}

	if done != nil {
		done <- err
		close(done)
	}
}
