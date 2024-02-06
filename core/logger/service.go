package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/thoas/go-funk"
)

// LogHTTPRequest log internal HTTP request
func LogHTTPRequest(ctx context.Context, host, url, method string, params map[string]interface{}) {
	requestStr := ""
	if params != nil {
		logParams := make(map[string]interface{})
		for k, v := range params {
			logParams[k] = v
		}

		marshalParams, _ := json.Marshal(logParams)
		requestStr = string(Secure(marshalParams))
	}

	logFields := map[string]interface{}{
		"host":   host,
		"url":    url,
		"method": method,
	}

	newURL, newURI := ParseURLAndURI(url)
	if !strings.EqualFold(newURI, "") {
		logFields["url"] = newURL
		logFields["uri"] = newURI
	}

	LogHTTPClientRequest(ctx, requestStr, logFields)
}

func ParseURLAndURI(url string) (newURL string, newURI string) {
	newURL = url
	if strings.Contains(url, "?") {
		urlArr := strings.Split(url, "?")
		if len(urlArr) > 1 {
			newURL = strings.ReplaceAll(url, urlArr[len(urlArr)-1], "{params}")
			newURI = url
		}
	} else {
		urlArr := strings.Split(url, "/")
		_, err := strconv.ParseInt(urlArr[len(urlArr)-1], 10, 64)
		if err == nil {
			newURL = strings.ReplaceAll(url, urlArr[len(urlArr)-1], "{id}")
			newURI = url
		} else {
			numeric := regexp.MustCompile(`\d`).MatchString(urlArr[len(urlArr)-1])
			if numeric {
				newURL = strings.ReplaceAll(url, urlArr[len(urlArr)-1], "{code}")
				newURI = url
			}
		}
	}
	return newURL, newURI
}

func Secure(data []byte) []byte {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(data, &bodyMap); err == nil {
		for k := range bodyMap {
			// Masking field in secretFields
			if funk.Contains(secretFields, k) {
				if _, ok := bodyMap[k]; ok {
					bodyMap[k] = "********"
					continue
				}
			}

			// Masking card number
			if k == "card_number" {
				if cardNumber, ok := bodyMap[k].(string); ok && len(cardNumber) > 12 {
					bodyMap[k] = fmt.Sprintf("%v%v%v", cardNumber[:6], "xxxxxx", cardNumber[12:])
					continue
				}
			}

		}

		maskData, _ := json.Marshal(bodyMap)
		return maskData
	}
	return data
}

// LogHTTPResponse log internal HTTP response
func LogHTTPResponse(ctx context.Context, host, url string, bodyResponse []byte, statusCode int, responseIn time.Duration) {
	responseStr := string(Secure(bodyResponse))

	logFields := map[string]interface{}{
		"host":          host,
		"url":           url,
		"status_code":   statusCode,
		"response_time": responseIn.String(),
	}

	newURL, newURI := ParseURLAndURI(url)
	if !strings.EqualFold(newURI, "") {
		logFields["url"] = newURL
		logFields["uri"] = newURI
	}

	LogHTTPClientResponse(ctx, responseStr, logFields)
}
