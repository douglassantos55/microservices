package pkg

import (
	"errors"
	"os"
	"reflect"
	"time"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	// Validates credencials and authenticates user
	Login(user, pass string) (string, string, error)

	// Validates and verifies token
	Verify(token string) (*User, error)
}

type service struct {
	tokenGen TokenGenerator
}

func NewService(tg TokenGenerator) Service {
	return &service{tg}
}

func (s *service) Login(user, pass string) (string, string, error) {
	if user != "admin" || pass != "123" {
		return "", "", ErrInvalidCredentials
	}

	payload := Payload{
		"iss": "auth",
		"exp": time.Now().Add(time.Hour),
		"sub": "admin",
		"aud": "renting",
		"user": &User{
			ID:   "aK0o3",
			Name: "John Doe",
		},
	}

	token, err := s.getToken(payload)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.getRefreshToken(payload)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (s *service) Verify(token string) (*User, error) {
	if token != "arandomtokenhere" {
		return nil, ErrInvalidToken
	}
	return &User{ID: "aK0o3", Name: "John Doe"}, nil
func (s *service) getToken(payload Payload) (string, error) {
	return s.tokenGen.Sign(payload, os.Getenv(JWT_SIGN_SECRET_ENV))
}

func (s *service) getRefreshToken(payload Payload) (string, error) {
	// refresh token expires only after a year
	payload["exp"] = time.Now().AddDate(1, 0, 0)
	return s.tokenGen.Sign(payload, os.Getenv(JWT_REFRESH_SECRET_ENV))
}
