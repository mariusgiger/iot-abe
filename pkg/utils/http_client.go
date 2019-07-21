package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/pkg/errors"
)

var httpTimeout = 20 * time.Second

// HTTPClient wraps an http.Client with optional tracing instrumentation.
type HTTPClient struct {
	Tracer opentracing.Tracer
	Client *http.Client
}

// NewHTTPClient creates a new tracing.HTTPClient
func NewHTTPClient(tracer opentracing.Tracer) *HTTPClient {
	return &HTTPClient{Client: &http.Client{Transport: &nethttp.Transport{}, Timeout: httpTimeout}, Tracer: tracer}
}

// Request type is used to compose and send individual request from client
type Request struct {
	client           *HTTPClient
	headers          map[string]string
	ctx              context.Context
	cancelFunc       context.CancelFunc
	jsonRequestBody  interface{}
	jsonResponseBody interface{}
}

// Response is an object represents executed request and its values.
type Response struct {
	Status  int
	RawBody []byte
}

// NewRequest creates a new request object, it is used form a HTTP/RESTful request
// such as GET, POST, PUT and DELETE
func (c *HTTPClient) NewRequest(ctx context.Context) *Request {
	return &Request{
		client: c,
		ctx:    ctx,
		cancelFunc: func() {
			//only used for timeouts (see WithTimeout)
		},
		headers: make(map[string]string),
	}
}

// WithHeaders method sets multiple header fields and their values in the current request.
// Example: To set `x-api-key` and `Accept` as `application/json`
//
// 		client.NewRequest(ctx).
//			WithHeaders(map[string]string{
//				"x-api-key": "someApiKey",
//				"Accept": "application/json",
//			})
//
func (r *Request) WithHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.headers[h] = v
	}
	return r
}

// WithHeader method is to set a single header field and its value in the current request.
// Example: To set `x-api-key` and `Accept` as `application/json`.
// 		resty.NewRequest(ctx).
//			WithHeader("x-api-key", "someApiKey").
//			WithHeader("Accept", "application/json")
//
func (r *Request) WithHeader(header, value string) *Request {
	r.headers[header] = value
	return r
}

// WithJSON method sets the request body for the request which is automatically marshalled to JSON
func (r *Request) WithJSON(body interface{}) *Request {
	r.jsonRequestBody = body
	return r
}

// WithResult method registers the response `Result` object for automatic unmarshalling
func (r *Request) WithResult(response interface{}) *Request {
	r.jsonResponseBody = response
	return r
}

// WithTimeout method performs the HTTP request with given timeout
// for current `Request`.
// 		client.NewRequest(ctx).WithTimeout(5000 * time.Millisecond)
//
func (r *Request) WithTimeout(timeout time.Duration) *Request {
	deadline := time.Now().Add(timeout)
	r.ctx, r.cancelFunc = context.WithDeadline(r.ctx, deadline)
	return r
}

// Get method executes a GET HTTP request
// 		resp, err := client.NewRequest(ctx).Get("http://gateway.sbrs.com")
//
func (r *Request) Get(url string) (*Response, error) {
	return r.Execute(http.MethodGet, url)
}

// Post method executes a POST HTTP request
// 		resp, err := client.NewRequest(ctx).Post("http://gateway.sbrs.com")
//
func (r *Request) Post(url string) (*Response, error) {
	return r.Execute(http.MethodPost, url)
}

// Put method executes a PUT HTTP request
// 		resp, err := client.NewRequest(ctx).Put("http://gateway.sbrs.com")
//
func (r *Request) Put(url string) (*Response, error) {
	return r.Execute(http.MethodPut, url)
}

// Delete method executes a DELETE HTTP request
// 		resp, err := client.NewRequest(ctx).Delete("http://gateway.sbrs.com")
//
func (r *Request) Delete(url string) (*Response, error) {
	return r.Execute(http.MethodDelete, url)
}

// Execute method performs the HTTP request with given HTTP method and URL
// for current `Request`.
// 		resp, err := client.NewRequest(ctx).Execute(http.MethodGet, "http://gateway.sbrs.com")
//
func (r *Request) Execute(method, url string) (*Response, error) {
	defer r.cancelFunc() //cancel request

	//prepare body
	var requestBody io.Reader
	if r.jsonRequestBody != nil {
		data, err := json.Marshal(r.jsonRequestBody)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request")
		}

		requestBody = bytes.NewBuffer(data)
	}

	//prepare request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare request")
	}

	//set headers
	for h, v := range r.headers {
		req.Header.Set(h, v)
	}
	req.Header.Set("Content-Type", "application/json")

	//set context and tracing
	req = req.WithContext(r.ctx)
	if r.client.Tracer != nil {
		var ht *nethttp.Tracer
		req, ht = nethttp.TraceRequest(r.client.Tracer, req, nethttp.OperationName(method+url)) //TODO does this produce useful logs?
		defer ht.Finish()
	}

	//do request
	res, err := r.client.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}
	defer MustClose(res.Body)

	//read response
	response := &Response{
		Status: res.StatusCode,
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, errors.Wrap(err, "could not read body")
	}
	response.RawBody = body

	//parse response
	if r.jsonResponseBody != nil {
		err = json.Unmarshal(body, r.jsonResponseBody)
		if err != nil {
			return response, errors.Wrapf(err, "could not decode body: %v", string(body))
		}
	}

	return response, nil
}
