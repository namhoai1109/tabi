package s2s

import "github.com/namhoai1109/tabi/core/server"

// ErrorResponseData model
type ErrorResponseData struct {
	Error server.HTTPError `json:"error"`
}

// InvokeOutputPayloadData model
type InvokeOutputPayloadData struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}
