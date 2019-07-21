package stream

import (
	"bufio"
	"bytes"
	"io"
)

var startCode = []byte{0xff, 0xd8}
var endCode = []byte{0xff, 0xd9}

// SplitMJPEGFrames splits an mjpeg stream
func SplitMJPEGFrames(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < len(startCode) {
		return 0, nil, nil
	}

	startFrame := bytes.Index(data[0:], startCode)
	if startFrame == -1 {
		return 0, nil, nil
	}

	endFrame := bytes.Index(data[startFrame:], endCode)
	if endFrame == -1 {
		return 0, nil, nil
	}

	frame := data[startFrame : endFrame+2]
	return endFrame + 2, frame, nil
}

// ReadImages reads an mjpeg stream from an io.Reader
//https://stackoverflow.com/questions/21702477/how-to-parse-mjpeg-http-stream-from-ip-camera
func ReadImages(src io.Reader, dest chan []byte, done chan bool, errCh chan error) {
	scanner := bufio.NewScanner(src)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	scanner.Split(SplitMJPEGFrames)

	// Read images from the buffer.
	for scanner.Scan() {
		jpg := scanner.Bytes()
		dest <- jpg
	}

	err := scanner.Err()
	if err != nil {
		errCh <- err
	}

	done <- true
}
