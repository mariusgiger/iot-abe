package stream

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//https://github.com/mattn/go-mjpeg/blob/master/mjpeg.go
//https://github.com/speps/grumpy-pi-mjpg
//https://github.com/hybridgroup/mjpeg

// Decoder decode motion jpeg
type Decoder struct {
	r *multipart.Reader
	m *sync.Mutex
}

// NewDecoder return new instance of Decoder
func NewDecoder(r io.Reader, b string) *Decoder {
	return &Decoder{
		r: multipart.NewReader(r, b),
		m: new(sync.Mutex),
	}
}

// NewDecoderFromResponse return new instance of Decoder from http.Response
func NewDecoderFromResponse(res *http.Response) (*Decoder, error) {
	_, param, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return NewDecoder(res.Body, strings.Trim(param["boundary"], "-")), nil
}

// NewDecoderFromURL return new instance of Decoder from response which specified URL
func NewDecoderFromURL(u string) (*Decoder, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return NewDecoderFromResponse(res)
}

// Decode do decoding
func (d *Decoder) Decode() (image.Image, error) {
	p, err := d.r.NextPart()
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(p)
}

// ProcessFunc used for processing stream
type ProcessFunc func(frame []byte) ([]byte, error)

// Stream wraps a streamer with different sinks
type Stream struct {
	m        *sync.Mutex
	s        map[chan []byte]struct{}
	Interval time.Duration
	log      *logrus.Logger
	process  ProcessFunc
}

// NewStream creates a new Stream
func NewStream(log *logrus.Logger) *Stream {
	return &Stream{
		m:   new(sync.Mutex),
		s:   make(map[chan []byte]struct{}),
		log: log,
	}
}

// WithProcessing sets a stream processing func
func (s *Stream) WithProcessing(p ProcessFunc) *Stream {
	s.process = p
	return s
}

// NewStreamWithInterval creates a new Stream which only streams in the specified interval and else drops frames (resource optimized)
func NewStreamWithInterval(interval time.Duration, log *logrus.Logger) *Stream {
	return &Stream{
		m:        new(sync.Mutex),
		s:        make(map[chan []byte]struct{}),
		Interval: interval,
		log:      log,
	}
}

// Close implements io.Closer
func (s *Stream) Close() error {
	s.m.Lock()
	defer s.m.Unlock()
	for c := range s.s {
		close(c)
		delete(s.s, c)
	}
	s.s = nil
	return nil
}

// Run starts streaming
func (s *Stream) Run(ctx context.Context, wg *sync.WaitGroup, camera *Camera) error {
	defer wg.Done()

	out := make(chan []byte)
	errChan := make(chan error)
	done := make(chan bool)

	s.log.Info("starting camera")
	err := camera.Start()
	if err != nil {
		return errors.Wrap(err, "could not start camera")
	}

	go camera.Read(out, done, errChan)

	for len(ctx.Done()) == 0 {
		var buf []byte

		select {
		case jpg := <-out:
			buf = jpg
		case err = <-errChan:
			s.log.Errorf("error while reading from camera: %v", err)
			continue
		case <-done:
			return nil
		}

		err = s.Update(buf)
		if err != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// Update updates all sinks
func (s *Stream) Update(b []byte) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.s == nil {
		return errors.New("stream was closed")
	}

	for c := range s.s {
		// Select to skip streams which are sleeping to drop frames.
		// This might need more thought.
		select {
		case c <- b:
		default:
		}
	}

	return nil
}

func (s *Stream) add(c chan []byte) {
	if s == nil {
		s.log.Error("pointer to s is nil in add")
	}

	if s.m == nil {
		s.log.Error("mutex is nil in add")
	}

	s.m.Lock()
	s.s[c] = struct{}{}
	s.m.Unlock()
}

func (s *Stream) destroy(c chan []byte) {
	s.m.Lock()
	if s.s != nil {
		close(c)
		delete(s.s, c)
	}
	s.m.Unlock()
}

// NWatch returns the number of watchers
func (s *Stream) NWatch() int {
	return len(s.s)
}

// Current returns the current frame
func (s *Stream) Current() []byte {
	c := make(chan []byte)
	s.add(c)
	defer s.destroy(c)

	return <-c
}

// WriteFiles sink used for testing
func (s *Stream) WriteFiles() {
	c := make(chan []byte)
	s.add(c)
	defer s.destroy(c)

	i := 0
	for {
		//time.Sleep(s.Interval)

		b, ok := <-c
		if !ok {
			fmt.Print("an error occurred while reading")
			break
		}

		fo, err := os.Create(fmt.Sprintf("img%v.jpg", i))
		if err != nil {
			fmt.Print("an error occurred while creating file")
			break
		}

		_, err = fo.Write(b)
		if err != nil {
			fmt.Print("an error occurred while writing file")
			break
		}

		i++
		if i >= 5 {
			return
		}
	}
}

// ServeHTTP serves the stream over http
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := make(chan []byte)
	s.add(c)
	defer s.destroy(c)

	m := multipart.NewWriter(w)
	defer m.Close()

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+m.Boundary())
	w.Header().Set("Connection", "close")
	h := textproto.MIMEHeader{}
	st := fmt.Sprint(time.Now().Unix())
	for {
		time.Sleep(s.Interval)

		b, ok := <-c
		if !ok {
			break
		}

		//NOTE this should be moved to stream.Update for better performance
		if s.process != nil {
			var err error
			b, err = s.process(b)
			if err != nil {
				resp := struct {
					Error string `json:"error"`
				}{
					Error: err.Error(),
				}
				utils.MustWriteJSON(w, resp, http.StatusInternalServerError)
			}
		}

		h.Set("Content-Type", "image/jpeg")
		h.Set("Content-Length", fmt.Sprint(len(b)))
		h.Set("X-StartTime", st)
		h.Set("X-TimeStamp", fmt.Sprint(time.Now().Unix()))
		mw, err := m.CreatePart(h)
		if err != nil {
			break
		}
		_, err = mw.Write(b)
		if err != nil {
			break
		}
		if flusher, ok := mw.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}
