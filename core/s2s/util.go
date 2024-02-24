package s2s

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	resty "github.com/go-resty/resty/v2"
	client "github.com/namhoai1109/tabi/core/http"
	"github.com/namhoai1109/tabi/core/logger"
	structutil "github.com/namhoai1109/tabi/util/struct"
	"github.com/thoas/go-funk"
)

// Get wrapper of request of method GET
func (s *Service) Get(ctx context.Context, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodGet, nil)
	resp, err := client.R().SetContext(ctx).Get(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// Post wrapper of request of method POST
func (s *Service) Post(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodPost, params)
	resp, err := client.R().SetContext(ctx).SetBody(params).Post(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// PostFormData wrapper of request of method POST with form-data
func (s *Service) PostFormData(ctx context.Context, params map[string]string, buff *bytes.Buffer, formName, fileName, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodPost, nil)
	resp, err := client.R().SetContext(ctx).SetFileReader(formName, fileName, buff).SetFormData(params).Post(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// Patch wrapper of request of method PATCH
func (s *Service) Patch(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodPatch, params)
	resp, err := client.R().SetContext(ctx).SetBody(params).Patch(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// PatchFormData wrapper of request of method PATCH with form-data
func (s *Service) PatchFormData(ctx context.Context, params map[string]string, buff *bytes.Buffer, formName, fileName, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodPost, nil)
	resp, err := client.R().SetContext(ctx).SetFileReader(formName, fileName, buff).SetFormData(params).Patch(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// Put wrapper of request of method PUT
func (s *Service) Put(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodPut, params)
	resp, err := client.R().SetContext(ctx).SetBody(params).Put(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// Delete wrapper of request of method DELETE
func (s *Service) Delete(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error) {
	client := client.NewClient(s.jwt)
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)
	client.GenerateAccessToken(time.Duration(s.cfg.Duration), customAccessToken...)
	logger.LogHTTPRequest(ctx, url, path, resty.MethodDelete, params)
	resp, err := client.R().SetContext(ctx).SetBody(params).Delete(path)
	if err != nil {
		return nil, err
	}
	logger.LogHTTPResponse(ctx, url, path, resp.Body(), resp.StatusCode(), resp.Time())
	return resp, nil
}

// BuildError to response error
func (s *Service) BuildError(resp *resty.Response) error {
	if resp.StatusCode() < http.StatusOK || resp.StatusCode() > http.StatusIMUsed {
		errResp := new(ErrorResponseData)
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			return err
		}
		return &errResp.Error
	}
	return nil
}

// InvokeLambda represents function to invoke lambda
func (s *Service) InvokeLambda(ctx context.Context, functionName, path, method string,
	headers map[string]string, body map[string]interface{}) (*InvokeLambdaResponse, error) {
	if s.cfg.Timeout == 0 {
		s.cfg.Timeout = 30
	}
	// Create Lambda client
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	lambdaClient := lambda.New(awsSession, &aws.Config{
		Region: aws.String(s.cfg.Region),
		HTTPClient: &http.Client{
			Timeout: time.Duration(time.Duration(s.cfg.Timeout).Seconds()),
		},
	})

	token := s.generateBasicToken(time.Duration(s.cfg.Duration))

	// Add headers
	reqHeaders := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %v", token),
	}
	for key, value := range headers {
		reqHeaders[key] = value
	}

	// Invoke Lambda function name with context
	logger.LogHTTPRequest(ctx, functionName, path, method, body)
	timeStart := time.Now()
	result, err := lambdaClient.InvokeWithContext(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      s.CreateInputPayLoad(reqHeaders, body, method, path),
	})
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("failed to invoke lambda function %v with err: %v", functionName, err.Error()))
		return nil, err
	}

	// Parse response and check status code
	resp := new(InvokeLambdaResponse)
	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		logger.LogError(context.Background(), fmt.Sprintf("failed to unmarshal response with err %v", err))
		return nil, err
	}
	resp.Body = strings.ReplaceAll(resp.Body, "\\n", "")

	logger.LogHTTPResponse(ctx, functionName, path, []byte(resp.Body), resp.StatusCode, time.Since(timeStart))

	if err := s.BuildErrorInvokeLambda(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// BuildErrorInvokeLambda to response error by invoke lambda
func (s *Service) BuildErrorInvokeLambda(errResp *InvokeLambdaResponse) error {
	if errResp.StatusCode < http.StatusOK || errResp.StatusCode > http.StatusIMUsed {
		body := new(ErrorResponseData)
		err := json.Unmarshal([]byte(errResp.Body), &body)
		if err != nil {
			logger.LogError(context.Background(), fmt.Sprintf("failed to unmarshal with err %v", err))
			return err
		}
		return &body.Error
	}

	return nil
}

func (s *Service) generateBasicToken(duration time.Duration) string {
	claims := map[string]interface{}{
		"role": client.RepresentativeRole, // by default, role is representative
	}
	expiredAt := time.Now().Add(time.Second * duration)
	token, _, err := s.jwt.GenerateToken(claims, &expiredAt)
	if err != nil {
		panic(err)
	}
	return token
}

func (s *Service) CreateInputPayLoad(reqHeaders map[string]string, body map[string]interface{}, method, path string) []byte {
	payload := make(map[string]interface{})
	payload["headers"] = reqHeaders
	payload["httpMethod"] = method
	payload["path"] = path
	if body != nil {
		bodyMarshal, _ := json.Marshal(body)
		payload["body"] = string(bodyMarshal)
	}
	payloadMarshaled, _ := json.Marshal(structutil.ToMap(payload))
	return payloadMarshaled
}

// GenerateAccessTokenByRole to generate access token
func (s *Service) GenerateAccessTokenByRole(id int, role string) (*string, error) {
	if funk.Contains([]string{
		client.RepresentativeRole,
		client.BranchManagerRole,
		client.ClientRole,
	}, role) {
		return nil, fmt.Errorf("role %s is not supported", role)
	}
	claims := map[string]interface{}{
		"id":   id,
		"role": role,
	}
	expiredAt := time.Now().Add(time.Minute * time.Duration(s.cfg.Duration))
	token, _, err := s.jwt.GenerateToken(claims, &expiredAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
