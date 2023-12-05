package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// ParseTokenFromHeader parses token from Authorization header
func (j *Service) parseTokenFromHeader(c echo.Context) (*jwt.Token, error) {
	// Verify basic token
	username, password, okBasic := c.Request().BasicAuth()
	if okBasic {
		return j.parseBasicToken(c.Request().Context(), username, password)
	}

	// Verify JWT
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, fmt.Errorf("token not found")
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
		return nil, fmt.Errorf("token invalid")
	}

	return j.parseToken(parts[1])
}

// ParseToken parses token from string
func (j *Service) parseToken(input string) (*jwt.Token, error) {
	return jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
		if j.algo != token.Method {
			return nil, fmt.Errorf("token method mismatched")
		}
		return j.key, nil
	})
}

// ParseBasicToken return token with claim of Backend ID
func (j *Service) parseBasicToken(ctx context.Context, username, password string) (*jwt.Token, error) {
	if j.secretCache == nil {
		fmt.Println("secretCache is nil. cannot use basic token")
		return nil, fmt.Errorf("token invalid")
	}

	basicTokenStr, err := j.secretCache.GetSecretString(j.cfg.SecretIDBasicToken)
	if err != nil {
		return nil, err
	}
	basicToken := new(BasicTokenData)
	err = json.Unmarshal([]byte(basicTokenStr), &basicToken)
	if err != nil {
		return nil, err
	}

	if username != basicToken.UserName {
		fmt.Println("basic token invalid")
		return nil, fmt.Errorf("token invalid")
	}

	if password != basicToken.Password {
		fmt.Println("basic token invalid")
		return nil, fmt.Errorf("token invalid")
	}

	token := &jwt.Token{
		Valid: true,
		Claims: jwt.MapClaims{
			"role": j.cfg.Role,
		},
	}

	return token, nil
}
