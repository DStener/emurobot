package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	virtualCamera = "/dev/video10"
	moduleName    = "v4l2loopback"
)

func main() {
	if err := ensureVirtualCamera(); err != nil {
		log.Fatal(err)
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal("ffmpeg not found")
	}

	args := append([]string{"ffmpeg"}, os.Args[1:]...)

	log.Println("starting ffmpeg")

	if err := syscall.Exec(ffmpegPath, args, os.Environ()); err != nil {
		log.Fatal(err)
	}
}

func ensureVirtualCamera() error {
	if _, err := os.Stat(virtualCamera); err == nil {
		log.Println("/dev/video10 already exists")
		return nil
	}

	log.Println("/dev/video10 not found, creating virtual camera")

	cmd := exec.Command(
		"modprobe",
		moduleName,
		"video_nr=10",
		"card_label=RobotVirtualCamera",
		"exclusive_caps=1",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		if _, err := os.Stat(virtualCamera); err == nil {
			log.Println("/dev/video10 created")
			return nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		time.Sleep(300 * time.Millisecond)
	}

	return errors.New("/dev/video10 was not created")
}
