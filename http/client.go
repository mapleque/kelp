package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Client A http client hold host and timeout
type Client struct {
	host    string
	timeout time.Duration
}

// NewClient Create a http client with host, default timeout is 10s.
func NewClient(host string) *Client {
	return &Client{
		host:    host,
		timeout: 10 * time.Second,
	}
}

// Request Send a request to target url with method and body.
// The response body has been read out, while others are in resp.
func Request(url, method string, body []byte) (*http.Response, []byte, error) {
	client := NewClient(url)
	client.SetTimeout(10 * time.Second)
	req, err := client.BuildRequest("", method, body)
	if err != nil {
		return nil, nil, err
	}
	return client.Do(req)
}

// Request Send a request to client host and target path with method and body.
// The response body has been read out, while others are in resp.
func (this *Client) Request(path, method string, body []byte) (*http.Response, []byte, error) {
	req, err := this.BuildRequest(path, method, body)
	if err != nil {
		return nil, nil, err
	}
	return this.Do(req)
}

// BuildRequest return a http.Request for any other expand.
// Use http.NewRequest to create a request is worked, too.
func (this *Client) BuildRequest(path, method string, body []byte) (*http.Request, error) {
	var data io.Reader
	if body != nil {
		data = bytes.NewReader(body)
	}
	req, err := http.NewRequest(
		method,
		this.host+path,
		data,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Do Send a request by using http.Request.
// The response body has been read out to body, while others are in resp.
func (this *Client) Do(req *http.Request) (resp *http.Response, body []byte, err error) {
	client := &http.Client{
		Timeout: this.timeout,
	}
	resp, err = client.Do(req)
	if err != nil {
		return resp, nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

// SetTimeout Set the timeout value of request.
func (this *Client) SetTimeout(d time.Duration) {
	this.timeout = d
}

// KelpClient A http client for kelp server, extend Client.
type KelpClient struct {
	Client
	token string
}

// NewKelpClient Create a KelpClient.
// The param token will be put into http header Authorization, which server may required.
// If the token is empty, there will no Authorization header.
// Default timeout is 10s.
func NewKelpClient(host, token string) *KelpClient {
	client := &KelpClient{}
	client.host = host
	client.token = token
	client.timeout = 10 * time.Second
	return client
}

// RequestKelp Send a request to a kelp server. Default timeout is 10s.
// The param token will be put into http header Authorization, which server may required.
// If the token is empty, there will no Authorization header.
// The param in is same as the handler param in, which define in the server.
// The param out is a response data holder, must be a pointer.
// The param lastContext is for request chain trace, such as traceid or uuid.
// The response status is not nil when the server returns an error status.
func RequestKelp(url, token string, in interface{}, out interface{}, lastContext *Context) (*Status, error) {
	client := NewKelpClient(url, token)
	client.SetTimeout(10 * time.Second)
	return client.RequestKelp("", in, out, lastContext)
}

// RequestKelp Send a request to the kelp server with path by client created before.
// The param in is same as the handler param in, which define in the server.
// The param out is a response data holder, must be a pointer.
// The param lastContext is for request chain trace, such as traceid or uuid.
// The response status is not nil when the server returns an error status.
func (this *KelpClient) RequestKelp(path string, in interface{}, out interface{}, lastContext *Context) (*Status, error) {
	var body []byte
	if in != nil {
		body, _ = json.Marshal(in)
	}
	req, err := this.BuildRequest(path, "POST", body)
	if err != nil {
		return nil, err
	}

	if len(this.token) > 0 {
		req.Header.Set("Authorization", this.token)
	}

	if lastContext != nil {
		req.Header.Set("Kelp-Traceid", lastContext.Request.Header.Get("Kelp-Traceid"))
		req.Header.Set("uuid", lastContext.Request.Header.Get("uuid"))
	}

	_, responseBody, err := this.Do(req)
	if err != nil {
		return nil, err
	}

	type responseType struct {
		Status  int             `json:"status"`
		Message interface{}     `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	response := &responseType{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, err
	}
	if response.Status != 0 {
		return JsonStatus(response.Status, response.Message), nil
	}

	if out == nil {
		return nil, nil
	}

	raw, _ := response.Data.MarshalJSON()
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	return nil, nil
}
