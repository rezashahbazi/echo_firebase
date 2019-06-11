package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"net/http"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
)

type (
	DefaultFirebaseConfig struct {
		mutex  sync.RWMutex
		Uptime time.Time `json:"uptime"`
		auth   *auth.Client
	}
)

func NewFirebaseMiddleWare(firebaseService *FirebaseService) *DefaultFirebaseConfig {
	return &DefaultFirebaseConfig{
		Uptime: time.Now(),
		auth:   firebaseService.app,
	}
}

const bearer = "Bearer"

type AuthFunc func(*auth.Token) (bool, error)

func AuthorizationFromParam(req *http.Request) (string, error) {
	return req.URL.Query().Get("authorization"), nil
}

func AuthorizationFromHeader(req *http.Request) (string, error) {
	header := req.Header.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("Authorization header not found")
	}

	l := len(bearer)
	if len(header) > l+1 && header[:l] == bearer {
		return header[l+1:], nil
	}

	return "", fmt.Errorf("Authorization header format must be 'Bearer {token}'")
}

func AuthorizationFromRequest(req *http.Request) (string, error) {
	authorization, err := AuthorizationFromParam(req)
	if authorization == "" {
		authorization, err = AuthorizationFromHeader(req)
		if err != nil {
			return "", err
		}
	}
	return authorization, nil
}

// Process is the middleware function.
func (s *DefaultFirebaseConfig) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		r := c.Request()
		authorization, err := AuthorizationFromRequest(r)

		if err != nil {
			log.Printf("error: %v", err)
			return next(c)
		}
		ctx := context.Background()
		token, err := s.auth.VerifyIDToken(ctx, authorization)

		if err != nil {
			log.Printf("error: %v", err)
			return next(c)
		}
		//putting UID to Echo Context
		c.Set("UID", token.UID)

		return next(c)
	}
}
