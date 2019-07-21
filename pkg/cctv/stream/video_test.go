package stream

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type VideoTestSuite struct {
	suite.Suite
}

func (ts *VideoTestSuite) TestMJPEGSplit() {
	//arrange
	videoFile, err := os.Open("./video.mjpg")
	ts.Require().NoError(err)

	videoFileInfo, err := videoFile.Stat()
	ts.Require().NoError(err)

	actual := make([]byte, 0)
	expected := make([]byte, videoFileInfo.Size())
	bytesRead, err := videoFile.Read(expected)
	fmt.Printf("read %v bytes\n", bytesRead)
	videoFile.Close()

	videoFile, err = os.Open("./video.mjpg")
	ts.Require().NoError(err)

	input := bufio.NewReader(videoFile)
	out := make(chan []byte)
	errChan := make(chan error)
	done := make(chan bool)

	//act
	go ReadImages(input, out, done, errChan)

	//assert
loop:
	for {
		select {
		case jpg := <-out:
			actual = append(actual, jpg...)
		case err = <-errChan:
			ts.Require().NoError(err)
			break loop
		case <-done:
			break loop
		}
	}

	ts.Require().Equal(len(expected), len(actual))
	ts.Require().Equal(expected, actual)
}

func (ts *VideoTestSuite) TestStream() {
	videoFile, err := os.Open("./video.mjpg")
	ts.Require().NoError(err)
	defer videoFile.Close()

	strm := NewStreamWithInterval(time.Millisecond*10, logrus.New())

	go strm.WriteFiles()

	input := bufio.NewReader(videoFile)
	out := make(chan []byte)
	errChan := make(chan error)
	doneChan := make(chan bool)

	go ReadImages(input, out, doneChan, errChan)

loop:
	for {
		select {
		case jpg := <-out:
			strm.Update(jpg)
		case err = <-errChan:
			ts.Require().NoError(err)
		case <-doneChan:
			time.Sleep(100 * time.Millisecond)
			break loop
		}
	}
}

func TestVideoTestSuite(t *testing.T) {
	suite.Run(t, &VideoTestSuite{})
}
