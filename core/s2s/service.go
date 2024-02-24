package s2s

import (
	"bytes"
	"context"

	resty "github.com/go-resty/resty/v2"
	"github.com/namhoai1109/tabi/core/middleware/jwt"
)

// Config represents the configuration
type Config struct {
	JwtAlgorithm string
	JwtSecret    string
	JwtDuration  int
	Duration     int // in seconds
	Debug        bool
	Timeout      int // in seconds
	Region       string
}

// New creates new bnpl service
func New(cfg Config) *Service {
	return &Service{
		jwt: jwt.New(cfg.JwtAlgorithm, cfg.JwtSecret, cfg.JwtDuration),
		cfg: cfg,
	}
}

// Service represents the bnpl service
type Service struct {
	jwt *jwt.Service
	cfg Config
}

// Intf represents interface
type Intf interface {
	Get(ctx context.Context, url, path string, customAccessToken ...string) (*resty.Response, error)
	Post(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error)
	PostFormData(ctx context.Context, params map[string]string, buff *bytes.Buffer, formName, fileName, url, path string, customAccessToken ...string) (*resty.Response, error)
	Patch(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error)
	PatchFormData(ctx context.Context, params map[string]string, buff *bytes.Buffer, formName, fileName, url, path string, customAccessToken ...string) (*resty.Response, error)
	Put(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error)
	Delete(ctx context.Context, params map[string]interface{}, url, path string, customAccessToken ...string) (*resty.Response, error)
	BuildError(resp *resty.Response) error

	InvokeLambda(ctx context.Context, functionName, path, method string, headers map[string]string, body map[string]interface{}) (*InvokeLambdaResponse, error)
	GenerateAccessTokenByRole(id int, role string) (*string, error)
}
