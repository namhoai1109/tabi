package secure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/namhoai1109/tabi-core/core/logger"
	"github.com/namhoai1109/tabi-core/core/middleware/logadapter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config represents secure specific config
type Config struct {
	AllowOrigins []string
}

// Headers adds general security headers for basic security measures
func Headers() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: true,
		// ContentSecurityPolicy: "default-src 'self'",
	})
}

// CORS adds Cross-Origin Resource Sharing support
func CORS(cfg *Config) echo.MiddlewareFunc {
	allowOrigins := []string{"*"}
	if cfg != nil && cfg.AllowOrigins != nil {
		allowOrigins = cfg.AllowOrigins
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           86400,
	})
}

// BodyDump prints out the request body for debugging purpose
func BodyDump() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		contentType := c.Request().Header.Get("Content-Type")
		if strings.EqualFold(contentType, "") {
			contentType = "application/json"
		}

		if len(reqBody) > 0 && !strings.EqualFold(contentType, "") && !strings.Contains(contentType, "multipart/form-data") {
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(reqBody, &bodyMap); err == nil {
				bodyMap = maskingBody(bodyMap)
				// Log body request
				logReqBody, _ := json.Marshal(bodyMap)
				logger.LogRequest(c, string(logReqBody), map[string]interface{}{
					"url": c.Request().URL.Path,
				})
			}
		}

		if len(resBody) > 0 && !strings.EqualFold(contentType, "") {
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(resBody, &bodyMap); err == nil {
				bodyMap = maskingBody(bodyMap)
				// Log body response
				logResBody, _ := json.Marshal(bodyMap)
				logger.LogResponse(c, string(logResBody), map[string]interface{}{
					"url": c.Request().URL.Path,
				})
			}
		}
	})
}

func maskingBody(bodyMap map[string]interface{}) map[string]interface{} {
	secretFields := []string{"new_password", "old_password", "password", "password_confirm", "access_token", "refresh_token", "otp", "otp_pass", "cvv", "voucher_code"}
	templateFields := []string{"otp", "otp_pass"}

	for k, v := range bodyMap {
		// Marking large fields
		if field, ok := v.(string); ok && len(field) > logadapter.DefaultLargeFieldLength {
			bodyMap[k] = fmt.Sprintf("%v%v%v", field[:50], "******", field[len(field)-50:])
		}
	}

	// Masking field in secretFields
	for i := 0; i < len(secretFields); i++ {
		if _, ok := bodyMap[secretFields[i]]; ok {
			bodyMap[secretFields[i]] = "********"
		}
	}

	// Masking card number
	if _, ok := bodyMap["card_number"]; ok {
		bodyMap["card_number"] = fmt.Sprintf("%v%v%v", bodyMap["card_number"].(string)[:6], "xxxxxx", bodyMap["card_number"].(string)[12:])
	}

	// Masking field in templateFields
	if templateData, ok := bodyMap["template_data"]; ok {
		templateDataByte, err := json.Marshal(templateData)
		if err != nil {
			return bodyMap
		}

		var templateDataMap map[string]interface{}
		if err := json.Unmarshal(templateDataByte, &templateDataMap); err == nil {
			for i := 0; i < len(templateFields); i++ {
				if _, ok := templateDataMap[templateFields[i]]; ok {
					templateDataMap[templateFields[i]] = "******"
				}
			}
		}
		bodyMap["template_data"] = templateDataMap
	}
	return bodyMap
}
