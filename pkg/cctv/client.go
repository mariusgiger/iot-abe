package cctv

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/mux"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/crypto"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Client wraps a simple web client
type Client struct {
	serverURL  string
	contract   common.Address
	user       common.Address
	manager    *acc.Manager
	wm         wallet.Manager
	httpClient *utils.HTTPClient
	log        *logrus.Logger
	cache      *cache.Cache
}

//NewClient creates a new client
func NewClient(serverURL string, contract, user common.Address, accManager *acc.Manager, wm wallet.Manager, log *logrus.Logger) *Client {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute) //TODO make this configurable

	return &Client{
		serverURL:  serverURL,
		contract:   contract,
		user:       user,
		manager:    accManager,
		wm:         wm,
		httpClient: utils.NewHTTPClient(nil),
		log:        log,
		cache:      c,
	}
}

//GetData retrieves data from the server and decrypts it
func (c *Client) GetData() (string, error) {
	acl, err := c.manager.ACLByAddress(c.user, c.contract)
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve acl")
	}

	encryptedKeyBytes, err := hexutil.Decode(acl.EncryptedKey)
	if err != nil {
		return "", errors.Wrap(err, "could not decode encrypte key")
	}

	keyBytes, err := c.wm.DecryptMessage(c.user, encryptedKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "could not decrypt priv key")
	}

	var data struct {
		Cipher string `json:"cipher"`
		Key    string `json:"key"`
	}

	url, err := utils.ParseURL(c.serverURL)
	if err != nil {
		return "", errors.Wrap(err, "could not parse server url")
	}
	url = url.AddPath("encrypt")

	resp, err := c.httpClient.
		NewRequest(context.Background()).
		WithResult(&data).
		Get(url.String())
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve data")
	}

	if resp.Status != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %v, %v", resp.Status, string(resp.RawBody))
	}

	pubKey, err := c.manager.PubKey(c.contract) //TODO cache
	if err != nil {
		return "", errors.Wrapf(err, "could not retrieve pubKey from contract (%v)", c.contract)
	}
	pubKeyBytes, err := hexutil.Decode(pubKey)
	if err != nil {
		return "", errors.Wrap(err, "could not decode pubKey")
	}

	cipherBytes, err := hexutil.Decode(data.Cipher)
	if err != nil {
		return "", errors.Wrap(err, "could not decode cipher text")
	}

	aesKeyBytes, err := hexutil.Decode(data.Key)
	if err != nil {
		return "", errors.Wrap(err, "could not decode aes key")
	}

	clearText, err := crypto.Decrypt(pubKeyBytes, keyBytes, aesKeyBytes, cipherBytes)
	if err != nil {
		return "", errors.Wrap(err, "could not decrypt message")
	}

	return string(clearText), nil
}

//GetImage retrieves an image from the server
func (c *Client) GetImage() (string, error) {
	url, err := utils.ParseURL(c.serverURL)
	if err != nil {
		return "", errors.Wrap(err, "could not parse server url")
	}
	url = url.AddPath("capture")

	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	path := "/tmp/asdf.jpg"
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

//ServeImage retrieves an image from the server
func (c *Client) ServeImage(bind string) error {
	c.log.Infof("starting iot client on: %v", bind)
	handler := c.MakeClientHandler()

	server := &http.Server{
		Addr:    bind,
		Handler: handler,
		// Good practice: enforce timeouts for servers you create!
		// does not work for streaming
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		c.log.Errorf("failed to start client: %v", err.Error())
		return errors.Wrap(err, "failed to start client")
	}

	return nil
}

// MakeClientHandler mounts all of the service endpoints into an http.Handler.
func (c *Client) MakeClientHandler() http.Handler {
	root := mux.NewRouter()

	root.Use(c.logging)

	root.Path("/capture").
		HandlerFunc(c.captureHandler).
		Methods(http.MethodGet)

	root.Path("/stream").
		HandlerFunc(c.serveStream).
		Methods(http.MethodGet)

	root.Path("/").
		HandlerFunc(c.streamSiteHandler()).
		Methods(http.MethodGet)

	return root
}

// logging is a simple logging middleware
func (c *Client) logging(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {

		L := c.log.WithContext(r.Context()).WithFields(logrus.Fields{
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

func (c *Client) serveStream(w http.ResponseWriter, r *http.Request) {
	m := multipart.NewWriter(w)
	defer m.Close()

	var data struct {
		Cipher string `json:"cipher"`
		Key    string `json:"key"`
	}

	privKey, err := c.getPrivKey()
	if err != nil {
		c.log.Fatal(errors.Wrap(err, "could not get private key"))
	}

	pubKey, err := c.getPubKey()
	if err != nil {
		c.log.Fatal(errors.Wrap(err, "could not get public key"))
	}

	url, err := utils.ParseURL(c.serverURL)
	if err != nil {
		c.log.Fatal(errors.Wrap(err, "could not parse server url"))
	}
	url.AddPath("stream")

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		c.log.Fatal(errors.Wrap(err, "could not create new request"))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.log.Fatal(errors.Wrap(err, "could not execute request"))
	}
	if resp.StatusCode != http.StatusOK {
		c.log.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
	}

	reader := bufio.NewReader(resp.Body)

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+m.Boundary())
	w.Header().Set("Connection", "close")
	h := textproto.MIMEHeader{}
	st := fmt.Sprint(time.Now().Unix())
	for {
		var image []byte
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}

				c.log.Fatal(errors.Wrap(err, "could not read from body buffer"))
			}

			if len(line) == 0 {
				continue
			}

			if bytes.Index(line, []byte("{")) != -1 {
				err = json.Unmarshal(line, &data)
				if err != nil {
					if err == io.EOF {
						break
					}
					c.log.Fatal(errors.Wrap(err, "could not decode encrypted image"))
				}

				cipherBytes, err := hexutil.Decode(data.Cipher)
				if err != nil {
					c.log.Fatal(errors.Wrap(err, "could not decode cipher text"))
				}

				aesKeyBytes, err := hexutil.Decode(data.Key)
				if err != nil {
					c.log.Fatal(errors.Wrap(err, "could not decode aes key"))
				}

				c.log.Info("decrypting image")

				img, err := crypto.Decrypt(pubKey, privKey, aesKeyBytes, cipherBytes)
				if err != nil {
					c.log.Fatal(errors.Wrap(err, "could not decrypt message"))
				}

				image = img
				break
			}
		}

		h.Set("Content-Type", "image/jpeg")
		h.Set("Content-Length", fmt.Sprint(len(image)))
		h.Set("X-StartTime", st)
		h.Set("X-TimeStamp", fmt.Sprint(time.Now().Unix()))

		mw, err := m.CreatePart(h)
		if err != nil {
			break
		}
		_, err = mw.Write(image)
		if err != nil {
			break
		}

		if flusher, ok := mw.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

const (
	privKeyCacheKey = "abe.privKey"
)

func (c *Client) getPubKey() ([]byte, error) {
	if key, found := c.cache.Get(pubKeyCacheKey); found {
		return key.([]byte), nil
	}

	pubKey, err := c.manager.PubKey(c.contract)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve pubKey from contract (%v)", c.contract)
	}
	pubKeyBytes, err := hexutil.Decode(pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode pubKey")
	}

	c.cache.Set(pubKeyCacheKey, pubKeyBytes, cache.DefaultExpiration)
	return pubKeyBytes, nil
}

func (c *Client) getPrivKey() ([]byte, error) {
	if key, found := c.cache.Get(privKeyCacheKey); found {
		return key.([]byte), nil
	}

	acl, err := c.manager.ACLByAddress(c.user, c.contract)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve acl")
	}

	encryptedKeyBytes, err := hexutil.Decode(acl.EncryptedKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode encrypte key")
	}

	keyBytes, err := c.wm.DecryptMessage(c.user, encryptedKeyBytes)
	if err != nil {
		c.log.Infof("key: %v", string(encryptedKeyBytes))
		return nil, errors.Wrap(err, "could not decrypt priv key")
	}

	c.cache.Set(privKeyCacheKey, keyBytes, cache.DefaultExpiration)
	return keyBytes, nil
}

// captureHandler handles GET /capture endpoint
func (c *Client) captureHandler(w http.ResponseWriter, r *http.Request) {
	image, status, err := func() ([]byte, int, error) {
		start := time.Now()
		privKey, err := c.getPrivKey()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		utils.TimeTrack(start, "get priv key")

		start = time.Now()
		pubKey, err := c.getPubKey()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		utils.TimeTrack(start, "get pub key")

		start = time.Now()
		var data struct {
			Cipher string `json:"cipher"`
			Key    string `json:"key"`
		}

		url, err := utils.ParseURL(c.serverURL)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not parse server url")
		}
		url = url.AddPath("captureenc")

		resp, err := c.httpClient.
			NewRequest(context.Background()).
			WithResult(&data).
			Get(url.String())
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not retrieve data")
		}

		if resp.Status != http.StatusOK {
			return nil, http.StatusInternalServerError, fmt.Errorf("unexpected status: %v, %v", resp.Status, string(resp.RawBody))
		}

		cipherBytes, err := hexutil.Decode(data.Cipher)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not decode cipher text")
		}

		aesKeyBytes, err := hexutil.Decode(data.Key)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not decode aes key")
		}
		utils.TimeTrack(start, "get image")

		start = time.Now()
		image, err := crypto.Decrypt(pubKey, privKey, aesKeyBytes, cipherBytes)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "could not decrypt message")
		}
		utils.TimeTrack(start, "decrypt image")

		return image, http.StatusOK, nil
	}()

	if err != nil {
		resp := &ErrorResponse{
			Error: err.Error(),
		}
		utils.MustWriteJSON(w, resp, status)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeContent(w, r, "camera.jpg", time.Now(), bytes.NewReader(image))
}

// streamSiteHandler handles GET / endpoint
func (c *Client) streamSiteHandler() func(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("webpage").Parse(defaultTemplate)
	if err != nil {
		c.log.Fatalf("could not parse template: %v", err)
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
		//TODO load metadata from server
		data.Len = 1 //converter.Len()
		err = t.Execute(w, data)
	}
}
