package paypal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/paypal/model"
	structutil "github.com/namhoai1109/tabi/util/struct"
)

var (
	GENERATE_ACCESS_TOKEN_PATH = "/v1/oauth2/token"
	CREATE_ORDER_PATH          = "/v2/checkout/orders"

	ORDER_STATUS_CREATED   = "CREATED"
	ORDER_STATUS_COMPLETED = "COMPLETED"
)

func (s *Service) generatePaypalClient(ctx context.Context) (*resty.Client, error) {
	client := resty.New()
	client.SetDebug(s.cfg.Debug)
	client.SetBaseURL(s.baseURL)
	client.SetTimeout(time.Duration(s.cfg.Timeout) * time.Second)

	body := map[string]interface{}{
		"grant_type":                "client_credentials",
		"ignoreCache":               true,
		"return_authn_schemes":      true,
		"return_client_metadata":    true,
		"return_unconsented_scopes": true,
	}

	logger.LogHTTPRequest(ctx, s.baseURL, GENERATE_ACCESS_TOKEN_PATH, resty.MethodPost, body)
	resp, err := client.R().
		SetContext(ctx).
		SetBasicAuth(s.cfg.ClientID, s.cfg.ClientSecret).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(s.convertToFormData(body)).
		Post(GENERATE_ACCESS_TOKEN_PATH)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when request to generate access token: %v", err))
		return nil, err
	}
	logger.LogHTTPResponse(ctx, s.baseURL, GENERATE_ACCESS_TOKEN_PATH, resp.Body(), resp.StatusCode(), resp.Time())

	if errResp := s.BuildError(resp); errResp != nil {
		return nil, errResp
	}

	data := new(model.AccessTokenResponse)
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when unmarshal response: %v", err))
		return nil, err
	}

	client.SetAuthToken(data.AccessToken).SetHeader("Content-Type", "application/json")
	return client, nil
}

func (s *Service) CreateOrder(ctx context.Context, creation *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	client, err := s.generatePaypalClient(ctx)
	if err != nil {
		return nil, err
	}

	logger.LogHTTPRequest(ctx, s.baseURL, CREATE_ORDER_PATH, resty.MethodPost, structutil.ToMap(creation))
	resp, err := client.R().
		SetContext(ctx).
		SetBody(creation).
		Post(CREATE_ORDER_PATH)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when request to create order: %v", err))
		return nil, err
	}
	logger.LogHTTPResponse(ctx, s.baseURL, CREATE_ORDER_PATH, resp.Body(), resp.StatusCode(), resp.Time())

	if errResp := s.BuildError(resp); errResp != nil {
		return nil, errResp
	}

	data := new(model.CreateOrderResponse)
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when unmarshal response: %v", err))
		return nil, err
	}

	return data, nil
}

func (s *Service) CaptureOrder(ctx context.Context, orderID string) (*model.CaptureOrderResponse, error) {
	client, err := s.generatePaypalClient(ctx)
	if err != nil {
		return nil, err
	}

	capturePath := CREATE_ORDER_PATH + "/" + orderID + "/capture"
	logger.LogHTTPRequest(ctx, s.baseURL, capturePath, resty.MethodPost, nil)
	resp, err := client.R().
		SetContext(ctx).
		Post(capturePath)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when request to capture order: %v", err))
		return nil, err
	}
	logger.LogHTTPResponse(ctx, s.baseURL, capturePath, resp.Body(), resp.StatusCode(), resp.Time())

	if errResp := s.BuildError(resp); errResp != nil {
		return nil, errResp
	}

	data := new(model.CaptureOrderResponse)
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when unmarshal response: %v", err))
		return nil, err
	}

	return data, nil
}

func (s *Service) BuildError(resp *resty.Response) *model.ErrorResponse {
	if resp.StatusCode() < http.StatusOK || resp.StatusCode() > http.StatusIMUsed {
		errResp := new(model.ErrorResponse)
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			return errResp
		}
		return errResp
	}
	return nil
}
