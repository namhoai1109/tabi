package logger

import (
	"context"

	"github.com/namhoai1109/tabi/core/middleware/logadapter"

	"github.com/labstack/echo/v4"
)

// LogError for logging errors with context to log request_id and correlation_id
func LogError(c context.Context, content interface{}) {
	logadapter.LogWithContext(c, content, logadapter.LogTypeError)
}

// LogDebug for logging debug with context to log request_id and correlation_id
func LogDebug(c context.Context, content interface{}) {
	logadapter.LogWithContext(c, content, logadapter.LogTypeDebug)
}

// LogInfo for logging info with context to log request_id and correlation_id
func LogInfo(c context.Context, content interface{}) {
	logadapter.LogWithContext(c, content, logadapter.LogTypeInfo)
}

// LogWarn for logging warn with context to log request_id and correlation_id
func LogWarn(c context.Context, content interface{}) {
	logadapter.LogWithContext(c, content, logadapter.LogTypeWarn)
}

// LogRequest for logging request with echo context to log request_id
func LogRequest(c echo.Context, content interface{}, logFields map[string]interface{}) {
	logadapter.LogWithEchoContext(c, content, logadapter.LogTypeRequest, logFields)
}

// LogResponse for logging response with echo context to log request_id
func LogResponse(c echo.Context, content interface{}, logFields map[string]interface{}) {
	logadapter.LogWithEchoContext(c, content, logadapter.LogTypeResponse, logFields)
}

// LogErrorWithEchoContext for logging errors with echo context to log request_id
func LogErrorWithEchoContext(c echo.Context, content interface{}) {
	logadapter.LogWithEchoContext(c, content, logadapter.LogTypeError)
}
