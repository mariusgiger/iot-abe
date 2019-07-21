package image

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

// TODO think about directly using the c-libraries
// https://github.com/raspberrypi/userland/tree/master/host_applications/linux/apps/raspicam
const (
	// RaspiStillBin path of raspistill (e.g.  "/usr/bin/raspistill")
	RaspiStillBin = "raspistill"
)

// Camera manages camera interactions with the raspistill command
type Camera interface {
	Capture() ([]byte, error)
}

type cameraService struct {
	log *logrus.Logger
}

// NewCamera creates a new Camera
func NewCamera(log *logrus.Logger) Camera {
	return &cameraService{
		log: log,
	}
}

// Capture takes an image with default settings and returnes the jpeg encoded byte array
func (s *cameraService) Capture() ([]byte, error) {
	//TODO compare https://github.com/dhowden/raspicam/blob/master/raspicam.go
	width := 1600
	height := 1200
	cameraParams := map[string]interface{}{
		"--timeout":    1,
		"--brightness": 50,
		"--quality":    90,
	}

	return s.CaptureRaspiStill(width, height, cameraParams)
}

// CaptureRaspiStill captures an image with given width, height, and other parameters
// return the captured image's bytes in jpeg encoding
// refer to: https://www.raspberrypi.org/documentation/usage/camera/raspicam/raspistill.md
func (s *cameraService) CaptureRaspiStill(width, height int, cameraParams map[string]interface{}) ([]byte, error) {
	// command line arguments
	args := []string{
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
		"-o", "-", // output to stdout
	}
	for k, v := range cameraParams {
		args = append(args, k)
		if v != nil {
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	// execute command
	bytes, err := exec.Command(RaspiStillBin, args...).CombinedOutput()
	if err != nil {
		s.log.Errorf("*** Error running %s: %s\n", RaspiStillBin, string(bytes))
		return []byte{}, err
	}

	return bytes, nil
}
