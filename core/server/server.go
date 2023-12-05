package server

import (
	"context"
	"fmt"
	"github.com/tabi-core/core/middleware/logadapter"
	"github.com/tabi-core/core/middleware/secure"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

// Config represents server specific config
type Config struct {
	Stage        string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	AllowOrigins []string
}

// DefaultConfig for the API server
var DefaultConfig = Config{
	Stage:        "development",
	Port:         3000,
	ReadTimeout:  10,
	WriteTimeout: 5,
	AllowOrigins: []string{"*"},
}

var echoLambda *echoadapter.EchoLambda

func (c *Config) fillDefaults() {
	if c.Stage == "" {
		c.Stage = DefaultConfig.Stage
	}
	if c.Port == 0 {
		c.Port = DefaultConfig.Port
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = DefaultConfig.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultConfig.WriteTimeout
	}
	if c.AllowOrigins == nil && len(c.AllowOrigins) == 0 {
		c.AllowOrigins = DefaultConfig.AllowOrigins
	}
}

// New instance new Echo server
func New(cfg *Config, isReqLog bool) *echo.Echo {
	cfg.fillDefaults()
	e := echo.New()

	if isReqLog {
		e.Use(logadapter.NewEchoLoggerMiddleware())
		e.Logger = logadapter.NewEchoLogger()
	}
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{LogLevel: log.ERROR}),
		secure.Headers(), secure.CORS(&secure.Config{AllowOrigins: cfg.AllowOrigins}))
	e.Validator = NewValidator()
	e.HTTPErrorHandler = NewErrorHandler(e).Handle
	e.Binder = NewBinder()
	e.Use(secure.BodyDump())
	e.Logger.SetLevel(log.DEBUG)
	logadapter.SetLevel(logadapter.DebugLevel)

	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Minute
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Minute
	return e
}

// Start starts echo server
func Start(e *echo.Echo, isDevelopment bool) {
	// graceful shutdown for dev environment
	if isDevelopment {
		// Start server
		go func() {
			if err := e.StartServer(e.Server); err != nil {
				if err == http.ErrServerClosed {
					fmt.Println("shutting down the server")
				} else {
					fmt.Errorf("error shutting down the server: %v", err)
				}
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal, 1) // buffered channel
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	} else {
		// User echo adapter for Lambda
		echoLambda = echoadapter.New(e)
		lambda.Start(Handler)
	}
}

// Handler function to handle request, response through Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}
