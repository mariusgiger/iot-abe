package cctv

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/cctv/image"
	"github.com/mariusgiger/iot-abe/pkg/cctv/stream"
	"github.com/mariusgiger/iot-abe/pkg/crypto"

	"github.com/gorilla/mux"
	"github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Server wraps a simple web server
type Server struct {
	log        *logrus.Logger
	bind       string
	cfg        *core.Config
	manager    *acc.Manager
	contract   common.Address
	deviceAddr common.Address
	camera     image.Camera
	cache      *cache.Cache
	stream     *stream.Stream
}

const defaultTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>Live</h1>
		<p>{{.Len}} viewers when page loaded</p>
		<img src="{{.Stream}}"/>
	</body>
</html>`

// NewServer creates a new server
func NewServer(bind string, log *logrus.Logger, cfg *core.Config, manager *acc.Manager, contract, device common.Address, camera image.Camera) *Server {

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute) //TODO make this configurable

	return &Server{
		log:        log,
		bind:       bind,
		cfg:        cfg,
		manager:    manager,
		contract:   contract,
		deviceAddr: device,
		camera:     camera,
		cache:      c,
	}
}

// Run starts the server and blocks
func (s *Server) Run() error {
	s.log.Infof("starting iot server on: %v", s.bind)
	ctx, cancel := context.WithCancel(context.Background())

	interval := 200 * time.Millisecond
	encrypt := func(image []byte) ([]byte, error) {
		pubKey, err := s.getPubKey()
		if err != nil {
			return nil, err
		}

		policy, err := s.getPolicy()
		if err != nil {
			return nil, err
		}
		s.log.Infof("encrypting data for: %v with policy %v", s.deviceAddr.Hex(), policy)

		//TODO program exits when encrypt fails for invalid policy (e.g. empty)
		//TODO allow reusing the symmetrical key
		key, cipher, err := crypto.Encrypt(pubKey, policy.Policy, image)
		if err != nil {
			return nil, errors.Wrap(err, "could not encrypt")
		}

		//TODO to encode this as the image is hacky...find a better way
		data := struct {
			Cipher string `json:"cipher"`
			Key    string `json:"key"`
		}{
			Cipher: hexutil.Encode(cipher),
			Key:    hexutil.Encode(key),
		}

		ct, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshall data")
		}

		return ct, nil
	}

	strm := stream.NewStreamWithInterval(interval, s.log).WithProcessing(encrypt)
	camera := stream.NewCamera(s.log)
	defer camera.Close()
	s.stream = strm

	handler := s.MakeHandler()

	server := &http.Server{
		Addr:    s.bind,
		Handler: handler,
		// Good practice: enforce timeouts for servers you create!
		// TODO define proper values for the stream
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go strm.Run(ctx, &wg, camera)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Errorf("failed to start server: %v", err.Error())
		s.log.Info("closing stream")
		strm.Close()
		cancel()
		return errors.Wrap(err, "failed to start server")
	}

	s.log.Info("closing stream")
	strm.Close()
	cancel()
	wg.Wait()

	return nil // OK
}

// MakeHandler mounts all of the service endpoints into an http.Handler.
func (s *Server) MakeHandler() http.Handler {
	root := mux.NewRouter()
	root.MethodNotAllowedHandler = http.HandlerFunc(s.notAllowed)
	root.NotFoundHandler = http.HandlerFunc(s.notFound)
	//root.Use(s.recoverPanic)
	root.Use(s.logging)
	root.Use(s.timeout(10 * time.Second)) // for all handlers

	root.Path("/encrypt").
		HandlerFunc(s.encryptHandler).
		Methods(http.MethodGet)

	//NOTE this endpoint only serves as a test endpoint and exposes camera without ac
	root.Path("/capture").
		HandlerFunc(s.captureHandler).
		Methods(http.MethodGet)

	root.Path("/stream").
		HandlerFunc(s.stream.ServeHTTP).
		Methods(http.MethodGet)

	root.Path("/").
		HandlerFunc(s.streamSiteHandler()).
		Methods(http.MethodGet)

	root.Path("/captureenc").
		HandlerFunc(s.captureEncryptedHandler).
		Methods(http.MethodGet)

	// GET version
	root.Path("/version").
		HandlerFunc(s.getVersion).
		Methods(http.MethodGet)

	return root
}

// encryptHandler handles GET /encrypt endpoint
func (s *Server) encryptHandler(w http.ResponseWriter, r *http.Request) {
	// get response and HTTP status
	resp, status, err := func() (interface{}, int, error) {
		pubKey, err := s.manager.PubKey(s.contract) //TODO cache
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrapf(err, "could not retrieve pubKey from contract (%v)", s.contract)
		}
		pubKeyBytes, err := hexutil.Decode(pubKey)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not decode pubKey")
		}

		policy, err := s.manager.DevicePolicyByAddress(s.deviceAddr, s.contract)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrapf(err, "could not retrieve device policy for (%v)", s.deviceAddr)
		}

		s.log.Infof("encrypting data for: %v with policy %v", s.deviceAddr, policy)

		//TODO program exits when encrypt fails for invalid policy (e.g. empty)
		key, cipher, err := crypto.Encrypt(pubKeyBytes, policy.Policy, []byte("Hello World!"))
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not encrypt")
		}

		return struct {
			Cipher string `json:"cipher"`
			Key    string `json:"key"`
		}{
			Cipher: hexutil.Encode(cipher),
			Key:    hexutil.Encode(key),
		}, http.StatusOK, nil
	}()

	utils.MustWriteJSON(w, s.checkError(err, resp), status)
}

// captureHandler handles GET /capture endpoint
func (s *Server) captureHandler(w http.ResponseWriter, r *http.Request) {
	image, err := s.camera.Capture()
	if err != nil {
		resp := &ErrorResponse{
			Error: err.Error(),
		}
		utils.MustWriteJSON(w, resp, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeContent(w, r, "camera.jpg", time.Now(), bytes.NewReader(image))
}

const (
	pubKeyCacheKey = "abe.pubkey"
	policyCacheKey = "device.policy"
)

func (s *Server) getPubKey() ([]byte, error) {
	if key, found := s.cache.Get(pubKeyCacheKey); found {
		return key.([]byte), nil
	}

	pubKey, err := s.manager.PubKey(s.contract)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve pubKey from contract (%v)", s.contract)
	}
	pubKeyBytes, err := hexutil.Decode(pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode pubKey")
	}

	s.cache.Set(pubKeyCacheKey, pubKeyBytes, cache.DefaultExpiration)
	return pubKeyBytes, nil
}

func (s *Server) getPolicy() (*acc.DevicePolicyEntry, error) {
	if policy, found := s.cache.Get(policyCacheKey); found {
		return policy.(*acc.DevicePolicyEntry), nil
	}

	policy, err := s.manager.DevicePolicyByAddress(s.deviceAddr, s.contract)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve device policy for (%v)", s.deviceAddr)
	}

	s.cache.Set(policyCacheKey, policy, cache.DefaultExpiration)
	return policy, nil
}

// captureEncryptedHandler handles GET /captureenc endpoint
func (s *Server) captureEncryptedHandler(w http.ResponseWriter, r *http.Request) {
	// get response and HTTP status
	resp, status, err := func() (interface{}, int, error) {
		start := time.Now()
		image, err := s.camera.Capture()
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrapf(err, "could not capture image")
		}
		utils.TimeTrack(start, "capture")

		start = time.Now()
		pubKey, err := s.getPubKey()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		utils.TimeTrack(start, "get pub key")

		start = time.Now()
		policy, err := s.getPolicy()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		utils.TimeTrack(start, "get policy")
		s.log.Infof("encrypting data for: %v with policy %v", s.deviceAddr.Hex(), policy)

		start = time.Now()
		//TODO program exits when encrypt fails for invalid policy (e.g. empty)
		//TODO allow reusing the symmetrical key
		key, cipher, err := crypto.Encrypt(pubKey, policy.Policy, image)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not encrypt")
		}
		utils.TimeTrack(start, "encrypt")

		return struct {
			Cipher string `json:"cipher"`
			Key    string `json:"key"`
		}{
			Cipher: hexutil.Encode(cipher),
			Key:    hexutil.Encode(key),
		}, http.StatusOK, nil
	}()

	utils.MustWriteJSON(w, s.checkError(err, resp), status)
}

// streamSiteHandler handles GET / endpoint
func (s *Server) streamSiteHandler() func(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("webpage").Parse(defaultTemplate)
	if err != nil {
		s.log.Fatalf("could not parse template: %v", err)
	}

	data := struct {
		Title  string
		Len    int
		Stream string
	}{
		Title:  "MJPG Server",
		Stream: "/stream",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//data.Len = s.broadcast.Len()
		err = t.Execute(w, data)
	}
}

//ErrorResponse is returned on error
type ErrorResponse struct {
	Error string `json:"error"`
}

// checkError checks the error and builds appropriate response body
func (s *Server) checkError(err error, resp interface{}) interface{} {
	if err != nil {
		// build error response
		return &ErrorResponse{
			Error: err.Error(),
		}
	}

	return resp
}

// notAllowed handler to get 405 error in proper format
func (s *Server) notAllowed(w http.ResponseWriter, r *http.Request) {
	L := s.log.WithContext(r.Context())
	L.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL,
	}).Error("method not allowed")

	resp := &ErrorResponse{Error: "method not allowed"}
	utils.MustWriteJSON(w, resp, http.StatusMethodNotAllowed)
}

// notFound handler to get 404 error in proper format
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	L := s.log.WithContext(r.Context())
	L.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL,
	}).Error("not found")

	resp := &ErrorResponse{Error: "not found"}
	utils.MustWriteJSON(w, resp, http.StatusNotFound)
}

// use middleware to handle panics and send 500 response back to user
func (s *Server) recoverPanic(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// panic?
			p := recover()
			if p == nil {
				return // done, no panic
			}

			// error?
			err, ok := p.(error)
			if !ok {
				panic(p) // not error, re-raise
			}

			L := s.log.WithContext(r.Context())
			L.WithFields(logrus.Fields{
				"method": r.Method,
				"url":    r.URL,
			}).Errorf("unhandled panic! %v", err)

			resp := &ErrorResponse{Error: err.Error()}
			utils.MustWriteJSON(w, resp, http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(mw)
}

// logging is a simple logging middleware
func (s *Server) logging(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {

		L := s.log.WithContext(r.Context()).WithFields(logrus.Fields{
			"remote": r.RemoteAddr,
			"method": r.Method,
			"url":    r.URL,
		})

		L.Info("request received")
		defer func(begin time.Time) {
			L.WithField("took", time.Since(begin)).Info("request completed")
		}(time.Now())

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(mw)
}

// timeout creates new middleware function with custom context timeout
func (s *Server) timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// getVersion handles `GET /version` endpoint
func (s *Server) getVersion(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]interface{}{
		"buildTime": s.cfg.BuildTime,
		"version":   s.cfg.Version,
		"gitHash":   s.cfg.GitHash,
		//"APIs": []string{"v1"},
	}
	utils.MustWriteJSON(w, resp, http.StatusOK)
}
