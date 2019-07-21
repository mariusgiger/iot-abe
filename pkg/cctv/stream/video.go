package stream

import (
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

const (
	// RaspiVidBin path of raspivid (e.g.  "/usr/bin/raspivid")
	RaspiVidBin = "raspivid"
	bufferSize  = 4096
)

// Camera wraps the raspivid command
type Camera struct {
	log    *logrus.Logger
	vidSrc io.Reader
}

// NewCamera creates a new camera
func NewCamera(log *logrus.Logger) *Camera {
	return &Camera{log: log}
}

// Read starts reading from the camera, the Start command has to be invoked beforehand
func (c *Camera) Read(dest chan []byte, done chan bool, errCh chan error) {
	if c.vidSrc == nil {
		c.log.Fatal("camera is not started")
	}

	ReadImages(c.vidSrc, dest, done, errCh)
}

// Start starts the raspivid command
func (c *Camera) Start() error {
	// command line arguments
	args := []string{
		"-cd", "MJPEG", //codec  H264 or MJPEG
		"-t", "0", //Time (in ms) to capture for. If not specified, set to 5s. Zero to disable
		"-n",        //Do not display a preview window
		"-hf",       // Set horizontal flip
		"-w", "640", //Set image width <size>
		"-h", "480", //Set image height <size>
		"-o", "-", // output to stdout
	}
	cmd := exec.Command(RaspiVidBin, args...)
	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		c.log.Errorf("could not read from stdout: %v", err)
		return err
	}
	c.vidSrc = stdoutIn

	if err = cmd.Start(); err != nil {
		c.log.Errorf("could not start command: %v", err)
		return err
	}
	c.log.Info("started raspivid command")

	return nil
}

// Close implements io.Closer
func (c *Camera) Close() error {
	if c.vidSrc != nil {
		if closer, ok := c.vidSrc.(io.ReadCloser); ok {
			return closer.Close()
		}
	}

	return nil
}
