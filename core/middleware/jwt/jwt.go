package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/namhoai1109/tabi/core/server"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Key string

var (
	USER_INFO_KEY Key = "X-User-Info"
)

// BasicTokenData model
type BasicTokenData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// New generates new JWT service necessery for auth middleware
func New(algo, secret string, duration int) *Service {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}
	return &Service{
		algo:     signingMethod,
		key:      []byte(secret),
		duration: time.Duration(duration) * time.Second,
	}
}

// NewWithConfig generates new JWT service with config necessery for auth middleware
func NewWithConfig(algo, secret string, duration int, config JWTConfig) *Service {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}
	secretCache, err := secretcache.New()
	if err != nil {
		fmt.Errorf("failed to new secretCache with err: %v", err)
	}
	return &Service{
		algo:        signingMethod,
		key:         []byte(secret),
		duration:    time.Duration(duration) * time.Second,
		cfg:         config,
		secretCache: secretCache,
	}
}

// JWTConfig represents config for JWT
type JWTConfig struct {
	SecretIDBasicToken string
	Role               string
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte
	// Duration (in seconds) for which the jwt token is valid.
	duration time.Duration
	// Service signing algorithm
	algo jwt.SigningMethod
	// Config
	cfg JWTConfig
	// Secret manager
	secretCache *secretcache.Cache
}

// MiddlewareFunction makes JWT implement the Middleware interface.
func (j *Service) MiddlewareFunction(services ...*Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.parseTokenFromHeader(c)
			if err != nil || !token.Valid {
				for _, svc := range services {
					if svc == nil {
						continue
					}
					t, svcErr := svc.parseTokenFromHeader(c)
					if svcErr != nil || !t.Valid {
						continue
					}
					claims := t.Claims.(jwt.MapClaims)
					info := make(map[string]interface{})
					for key, val := range claims {
						c.Set(key, val)
						info[key] = val
					}
					ctx := c.Request().Context()
					ctx = context.WithValue(ctx, USER_INFO_KEY, info)
					request := c.Request().WithContext(ctx)
					c.SetRequest(request)
					return next(c)
				}
				if err != nil {
					fmt.Errorf("error parsing token: %+v", err)
				}
				return server.NewHTTPAuthorizationError("Your session is unauthorized or has expired.")
			}
			claims := token.Claims.(jwt.MapClaims)
			info := make(map[string]interface{})
			for key, val := range claims {
				c.Set(key, val)
				info[key] = val
			}
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, USER_INFO_KEY, info)
			request := c.Request().WithContext(ctx)
			c.SetRequest(request)
			return next(c)
		}
	}
}

// GenerateToken generates new Service token and populates it with user data
func (j *Service) GenerateToken(claims map[string]interface{}, expire *time.Time) (string, int, error) {
	if expire == nil {
		expTime := time.Now().Add(j.duration)
		expire = &expTime
	}
	claims["exp"] = expire.Unix()

	token := jwt.NewWithClaims(j.algo, jwt.MapClaims(claims))
	tokenString, err := token.SignedString(j.key)

	return tokenString, int(time.Until(*expire).Seconds()), err
}
