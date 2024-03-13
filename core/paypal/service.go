package paypal

import (
	"context"

	"github.com/namhoai1109/tabi/core/paypal/model"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Debug        bool
	Timeout      int // in seconds
}

func New(baseURL string, cfg Config) *Service {
	return &Service{
		baseURL: baseURL,
		cfg:     cfg,
	}
}

type Service struct {
	baseURL string
	cfg     Config
}

type Intf interface {
	CreateOrder(ctx context.Context, creation *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	CaptureOrder(ctx context.Context, orderID string) (*model.CaptureOrderResponse, error)
	GetOrderDetails(ctx context.Context, orderID string) (*model.GetOrderDetailResponse, error)
}
