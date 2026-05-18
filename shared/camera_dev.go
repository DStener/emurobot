package emurobot

import (
	"errors"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	inputCamera   = "/dev/video0"
	virtualCamera = "/dev/video10"

	videoSize   = "640x480"
	fps         = "5"
	inputFormat = "yuyv422"
)



func ensureVirtualCamera() error {
	if _, err := os.Stat(virtualCamera); err == nil {
		log.Println("/dev/video10 already exists")
		return nil
	}

	log.Println("/dev/video10 not found, creating virtual camera")

	cmd := exec.Command(
		"modprobe",
		"v4l2loopback",
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

func InitCamera(config Config) error {
	if err := ensureVirtualCamera(); err != nil {
		return err
	}

	return nil
}