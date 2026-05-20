package main

import (
	"os"
	"os/exec"
	"sync"
	"syscall"

	emurobot "emurobot/shared"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var (
	mu           sync.Mutex
	videoFileCmd *exec.Cmd
)

func InitPlayer(path string) (string, error) {

	err := emurobot.LoadDumps(path)
	if err != nil {
		return "ERROR", err
	}

	for i := 0; i < emurobot.GetDeviceCount(); i++ {
		dump := emurobot.GetDeviceLog(i)
		go runGhostPlayer(*dump)
	}

	return "START", nil
}

func InitCameraPlayer(path string) (string, error) {

	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-stream_loop", "-1",
		"-re",
		"-i", path,
		"-vf", "scale=640:480,format=yuyv422",
		"-vcodec", "rawvideo",
		"-pix_fmt", "yuyv422",
		"-f", "v4l2",
		"/dev/video10",
	}

	cmd := exec.Command("ffmpeg", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return "ERROR", err
	}

	videoFileCmd = cmd

	go func() {
		err := cmd.Wait()

		mu.Lock()
		defer mu.Unlock()

		if videoFileCmd == cmd {
			videoFileCmd = nil
		}

		if err != nil {
			log.Println("video player stopped:", err)
		} else {
			log.Println("video player finished normally")
		}
	}()

	log.Println("video player started:", path)

	return "START", nil
}

func StopCameraPlayer() error {
	if videoFileCmd == nil {
		return nil
	}

	pgid, err := syscall.Getpgid(videoFileCmd.Process.Pid)
	if err == nil {
		err = syscall.Kill(-pgid, syscall.SIGINT)
	} else {
		err = videoFileCmd.Process.Kill()
	}

	videoFileCmd = nil

	return err
}

func runGhostPlayer(dump emurobot.LogDump) {

	// Wait until devise is not exist
	errIn := emurobot.WaitDeviceExist(dump.Input)
	errOut := emurobot.WaitDeviceExist(dump.Output)

	if errIn != nil || errOut != nil {
		log.Error(errIn, errOut)
	}

	// Create configs
	inputConfig := serial.Config{
		Name:        dump.Input,
		Baud:        dump.Speed,
		ReadTimeout: 1 * time.Second,
		Size:        byte(dump.Size),
	}

	input, err := serial.OpenPort(&inputConfig)
	if err != nil {
		log.Panicln("Not open input port", err.Error())
	}
	defer input.Close()

	// Loop for simulate real device
	for _, event := range dump.Events {
		// Simulate input time
		time.Sleep(time.Duration(event.Time))

		// Write to ghost port
		_, err := input.Write([]byte(event.Data))
		if err != nil {
			log.Fatal(err)
		}
	}
}
