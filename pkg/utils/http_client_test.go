package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
	tracer opentracing.Tracer
	closer io.Closer
}

func (suite *HTTPTestSuite) SetupSuite() {
	tracer, closer := InitJaeger("httptest")
	suite.tracer = tracer
	suite.closer = closer
}

type response struct {
	Msg string `json:"msg,omitempty"`
	Err string `json:"error,omitempty"`
}
type request struct {
	ReqMsg string `json:"reqMsg,omitempty"`
}

func (suite *HTTPTestSuite) TestGetJSON() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "someMsg",
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		WithResult(responseStruct).
		Get(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), "someMsg", responseStruct.Msg)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
}

func (suite *HTTPTestSuite) TestGetJSONWithHeaders() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), "someApiKey", r.Header.Get("x-api-key"))
		assert.Equal(suite.T(), "application/json", r.Header.Get("Accept"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "someMsg",
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		WithResult(responseStruct).
		WithHeaders(map[string]string{
			"x-api-key": "someApiKey",
			"Accept":    "application/json",
		}).
		Get(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), "someMsg", responseStruct.Msg)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
}

func (suite *HTTPTestSuite) TestGetJSONWithHeader() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), "someApiKey", r.Header.Get("x-api-key"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "someMsg",
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		WithHeader("x-api-key", "someApiKey").
		WithResult(responseStruct).
		Get(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), "someMsg", responseStruct.Msg)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
}

func (suite *HTTPTestSuite) TestGetJSONDeadlineExceeded() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 7)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "someMsg",
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		WithResult(responseStruct).
		WithTimeout(time.Second * 5).
		Get(remoteStub.URL)

	// assert
	require.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "context deadline exceeded")
	assert.Nil(suite.T(), resp)
}

func (suite *HTTPTestSuite) TestGetJSONErrorReturned() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "someErrorOccurred",
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		WithResult(responseStruct).
		WithTimeout(time.Second * 5).
		Get(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	require.NotEmpty(suite.T(), responseStruct.Err)
	assert.Contains(suite.T(), responseStruct.Err, "someErrorOccurred")
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Status)
}

func (suite *HTTPTestSuite) TestGetJSONNoResult() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "someMsg",
		})

	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		Get(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
	assert.NotEmpty(suite.T(), resp.RawBody)
}

func (suite *HTTPTestSuite) TestPostJSON() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		decoder := json.NewDecoder(r.Body)
		req := &request{}
		decoder.Decode(req)

		json.NewEncoder(w).Encode(map[string]string{
			"msg": req.ReqMsg,
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	data := &request{"someRequestMsg"}
	resp, err := client.
		NewRequest(context.Background()).
		WithJSON(data).
		WithResult(responseStruct).
		Post(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), "someRequestMsg", responseStruct.Msg)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
}

func (suite *HTTPTestSuite) TestPutJSON() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodPut, r.Method)
		w.WriteHeader(http.StatusOK)
		decoder := json.NewDecoder(r.Body)
		req := &request{}
		decoder.Decode(req)

		json.NewEncoder(w).Encode(map[string]string{
			"msg": req.ReqMsg,
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	data := &request{"someRequestMsg"}
	resp, err := client.
		NewRequest(context.Background()).
		WithJSON(data).
		WithResult(responseStruct).
		Put(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), "someRequestMsg", responseStruct.Msg)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
}

func (suite *HTTPTestSuite) TestDeleteJSON() {
	// arrange
	remoteStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "someMsg",
		})
	}))
	defer remoteStub.Close()

	client := &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}}, Tracer: suite.tracer}
	responseStruct := &response{}

	// act
	resp, err := client.
		NewRequest(context.Background()).
		WithResult(responseStruct).
		Delete(remoteStub.URL)

	// assert
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), "someMsg", responseStruct.Msg)
	assert.Equal(suite.T(), http.StatusOK, resp.Status)
}

func (suite *HTTPTestSuite) TearDownSuite() {
	suite.closer.Close()
}

func TestHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}
