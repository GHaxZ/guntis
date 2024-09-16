// Package jsonrpc - Code for performing json
package jsonrpc

import (
	"net/http"
	"net/http/cookiejar"
)

// Client - Jsonrpc client useful for working with a jsonrpc API
type Client struct {
	url       string
	client    *http.Client
	cookieJar *cookiejar.Jar
}

// Params - A collection of parameters used for a jsonrpc request
type Params struct {
	params map[string]interface{}
}

// Headers - A collection of headers used for a jsonrpc request
type Headers struct {
	headers map[string]string
}

// Request - Represents a jsonrpc request
type Request struct {
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	JSONRPC string      `json:"jsonrpc"`
}

// Response - Represents a jsonrpc response
type Response struct {
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error,omitempty"`
	JSONRPC string      `json:"jsonrpc"`
}

// Error - Represents a jsonrpc error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
