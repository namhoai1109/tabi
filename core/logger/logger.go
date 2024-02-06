package logger

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/namhoai1109/tabi/core/middleware/logadapter"
)

var secretFields = []string{"new_password", "old_password", "password", "password_confirm", "new_password_confirm", "access_token", "refresh_token", "otp", "otp_pass", "voucher_code"}

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

// LogHTTPClientResponse for logging HTTP Client response with correlation_id
func LogHTTPClientResponse(c context.Context, content interface{}, logFields map[string]interface{}) {
	logadapter.LogWithContext(c, content, logadapter.LogTypeHTTPClientResponse, logFields)
}

// LogHTTPClientRequest for logging HTTP Client request with correlation_id
func LogHTTPClientRequest(c context.Context, content interface{}, logFields map[string]interface{}) {
	logadapter.LogWithContext(c, content, logadapter.LogTypeHTTPClientRequest, logFields)
}
