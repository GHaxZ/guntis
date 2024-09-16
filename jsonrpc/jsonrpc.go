// Package jsonrpc - Code for performing jsonrpc communication
package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// NewClient - Create a new jsonrpc client
func NewClient(url string) (*Client, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, fmt.Errorf("Failed creating cookie jar for client: %s", err.Error())
	}

	return &Client{
		url:       url,
		client:    &http.Client{Jar: jar},
		cookieJar: jar,
	}, nil
}

// NewParams - Create a new collection of parameters
func NewParams() Params {
	return Params{
		params: make(map[string]interface{}, 0),
	}
}

// NewHeaders - Create a new collection of headers
func NewHeaders() Headers {
	return Headers{
		headers: make(map[string]string, 0),
	}
}

// SendRequest - Send a request using this client
func (c *Client) SendRequest(request Request) (*Response, error) {
	return c.SendRequestH(request, NewHeaders())
}

// Add - Add a parameter
func (p Params) Add(key string, value interface{}) Params {
	p.params[key] = value

	return p
}

// Add - Add a header
func (h Headers) Add(key, value string) Headers {
	h.headers[key] = value

	return h
}

// SendRequestH - Send a request using this client with headers
func (c *Client) SendRequestH(request Request, headers Headers) (*Response, error) {
	reqJSON, err := request.Serialize()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(reqJSON))
	req.Header.Add("Content-Type", "application/json")

	for key, value := range headers.headers {
		req.Header.Add(key, value)
	}

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// AddCookie - Add a cookie to the client manually
func (c *Client) AddCookie(cookie http.Cookie) error {
	u, err := url.Parse(c.url)

	if err != nil {
		return fmt.Errorf("Failed adding cookie: %s", err.Error())
	}

	c.cookieJar.SetCookies(u, []*http.Cookie{&cookie})

	return nil
}

// GetCookies - Get cookies for the current URL
func (c *Client) GetCookies() ([]*http.Cookie, error) {
	u, err := url.Parse(c.url)

	if err != nil {
		return nil, fmt.Errorf("Failed getting cookies: %s", err.Error())
	}

	return c.cookieJar.Cookies(u), nil
}

// NewRequest - Build a new jsonrpc request
func NewRequest(method string, params Params) Request {
	return Request{
		ID:      "0",
		Method:  method,
		Params:  params.params,
		JSONRPC: "2.0",
	}
}

// NewRequestRaw - Build a new jsonrpc request, which accepts raw parameters
func NewRequestRaw(method string, rawParams string) Request {
	return Request{
		ID:      "0",
		Method:  method,
		Params:  json.RawMessage(rawParams),
		JSONRPC: "2.0",
	}
}

// Serialize - Serialize a jsonrpc request
func (r *Request) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

// Format - Format this error into a readable string
func (e *Error) Format() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// Content - Get the content of the response, either a result or an error
func (r *Response) Content() (interface{}, *Error) {
	if r.Error != nil {
		return nil, r.Error
	}

	return r.Result, nil
}
