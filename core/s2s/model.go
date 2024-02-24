package s2s

import (
	"github.com/namhoai1109/tabi/core/server"
)

// ErrorResponseData model
type ErrorResponseData struct {
	Error server.HTTPError `json:"error"`
}

// InvokeLambdaResponse model
type InvokeLambdaResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}
