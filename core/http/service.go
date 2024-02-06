package httpcore

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	dbcore "github.com/namhoai1109/tabi/core/db"
	structutil "github.com/namhoai1109/tabi/util/struct"
	"gorm.io/gorm"
)

type key string
type Response = resty.Response
type Request = resty.Request

// Exported constanst
const (
	CorrelationIDKey key = "X-User-Correlation-Id"
)

type (
	// Client extends resty client
	Client[T any] struct {
		*resty.Client
		jwt     JWT
		db      *gorm.DB
		log     LogDB[T]
		sender  chan T
		close   chan bool
		wait    chan bool
		payload *T
	}

	// JWT interface
	JWT interface {
		GenerateToken(claims map[string]interface{}, expire *time.Time) (string, int, error)
	}

	// LogDB repository interface
	LogDB[T any] interface {
		dbcore.Intf
		SafeWrite(db *gorm.DB, data *T) (*T, error) // TODO: NEED TO CHECK THIS
	}
)

// NewClient creates new http client
func NewClient(jwt JWT) *Client[any] {
	return &Client[any]{
		jwt:    jwt,
		Client: resty.New(),
	}
}

// NewClientWithLog creates new http client
func NewClientWithLog[T any](db *gorm.DB, logDB LogDB[T], jwt JWT) *Client[T] {
	c := &Client[T]{
		jwt:    jwt,
		Client: resty.New(),
		db:     db,
		log:    logDB,
		sender: make(chan T),
		close:  make(chan bool),
		wait:   make(chan bool),
	}
	go c.handle()
	return c
}

// UseRequestCallBack process before client make request
func (c *Client[T]) UseRequestCallBack(f func(req *Request, payload *T) T) *Client[T] {
	c.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		c.sender <- f(request, c.payload)
		return nil
	})
	return c
}

// UseResponseCallBack process after http client response
func (c *Client[T]) UseResponseCallBack(f func(res *Response, payload *T) T) *Client[T] {
	c.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		//Prevent deadlock
		time.AfterFunc(3*time.Second, func() {
			c.wait <- true
		})
		<-c.wait
		c.sender <- f(response, c.payload)
		return nil
	})
	return c
}

// UseErrorCallBack process when http client returns error
func (c *Client[T]) UseErrorCallBack(f func(req *Request, payload *T, err *error) T) *Client[T] {
	c.OnError(func(request *resty.Request, err error) {
		//Prevent deadlock
		time.AfterFunc(3*time.Second, func() {
			c.wait <- true
		})
		<-c.wait
		c.sender <- f(request, c.payload, &err)
	})
	return c
}

// GenerateAccessToken generates access token
func (c *Client[any]) GenerateAccessToken(duration time.Duration, customAccessToken ...string) *Client[any] {
	var accessToken string
	if len(customAccessToken) > 0 {
		accessToken = customAccessToken[0]
	} else {
		claims := map[string]interface{}{
			"role": "admin", // by default, role is admin
		}
		expiredAt := time.Now().Add(time.Second * duration)
		token, _, err := c.jwt.GenerateToken(claims, &expiredAt)
		if err != nil {
			panic(err)
		}
		accessToken = token
	}
	c.SetHeaders(map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", accessToken)})
	return c
}

func (c *Client[T]) handle() {
	for {
		select {
		case data := <-c.sender:
			if c.payload == nil {
				r, err := c.log.SafeWrite(c.db, &data)
				if err != nil {
					fmt.Printf("Return error when insert purchase log, err:%v", err)
				} else {
					c.payload = r
					c.wait <- true
				}
			} else {
				// Only update existed purchase
				logMap := structutil.ToMap(data)
				if id, ok := logMap["id"]; ok {
					if err := c.log.Update(c.db, data, ` id  = ?`, id); err != nil {
						fmt.Printf("Return error when update log [%d], err:%v", id, err)
					}
					c.payload = &data
				} else {
					fmt.Println("Cannot get ID")
				}
			}
		case <-c.close:
			return
		}
	}
}

// CloseLog to close on related channels
func (c Client[T]) CloseLog() {
	c.close <- true
	close(c.close)
	close(c.sender)
}
