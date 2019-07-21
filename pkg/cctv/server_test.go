package cctv

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
	"gopkg.in/gavv/httpexpect.v1"
)

type ServerTestSuite struct {
	suite.Suite

	server     *Server
	testServer *httptest.Server
	expect     *httpexpect.Expect
}

// SetupSuite creates server instance
func (ts *ServerTestSuite) SetupSuite() {
	ts.server = &Server{log: logrus.New(), cfg: &core.Config{}}
	ts.testServer = httptest.NewServer(ts.server.MakeHandler())
	ts.expect = httpexpect.New(ts.T(), ts.testServer.URL)
}

// TestGetVersion integration tests for `GET /version`
func (ts *ServerTestSuite) TestGetVersion() {
	ts.server.cfg.GitHash = "5af68a50e445db0d8283beda71a5d279f04aa0c4"
	ts.server.cfg.BuildTime = "2019-04-12T19:47:05Z"
	ts.server.cfg.Version = "0.0.1"

	// missing URL
	obj := ts.expect.GET("/versionMissing").
		Expect().Status(http.StatusNotFound).JSON().Object()
	obj.Schema(&ErrorResponse{})
	obj.Path(`$.error`).String().
		Contains("not found")

	// not GET
	obj = ts.expect.POST("/version").
		Expect().Status(http.StatusMethodNotAllowed).JSON().Object()
	obj.Schema(&ErrorResponse{})
	obj.Path(`$.error`).String().
		Contains("method not allowed")

	// positive case
	obj = ts.expect.
		GET("/version").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Path(`$.gitHash`).String().Equal(ts.server.cfg.GitHash)
	obj.Path(`$.buildTime`).String().Equal(ts.server.cfg.BuildTime)
	obj.Path(`$.version`).String().Equal(ts.server.cfg.Version)
}

// TestGetImage integration tests for `GET /capture`
func (ts *ServerTestSuite) TestGetImage() {
	mock := &MockCameraService{}
	ts.server.camera = mock
	img := ts.mockImage()

	mock.resp = nil
	mock.err = errors.New("some error")
	obj := ts.expect.
		GET("/capture").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().
		Object()
	obj.Schema(&ErrorResponse{})
	obj.Path(`$.error`).String().
		Contains("some error")

	mock.resp = img
	mock.err = nil
	resp := ts.expect.
		GET("/capture").
		Expect().
		Status(http.StatusOK).
		ContentType("image/jpeg")
	body := []byte(resp.Body().Raw())
	ts.Equal(img, body)
}

func (ts *ServerTestSuite) mockImage() []byte {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	buffer := new(bytes.Buffer)

	err := jpeg.Encode(buffer, m, nil)
	ts.Require().NoError(err)

	return buffer.Bytes()
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}

type MockCameraService struct {
	err  error
	resp []byte
}

func (s *MockCameraService) Capture() ([]byte, error) {
	return s.resp, s.err
}
