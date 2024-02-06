package s2s

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	resty "github.com/go-resty/resty/v2"
	client "github.com/namhoai1109/tabi/core/http"
	"github.com/namhoai1109/tabi/core/logger"
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
		error := new(ErrorResponseData)
		json.Unmarshal(resp.Body(), &error)
		return &error.Error
	}
	return nil
}
