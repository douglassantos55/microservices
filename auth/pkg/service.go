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

type AuthResponse struct {
	User    *User  `json:"user"`
	Token   string `json:"token"`
	Refresh string `json:"refresh_token"`
}

type Service interface {
	// Validates credencials and authenticates user
	Login(user, pass string) (*AuthResponse, error)

	// Validates and verifies token
	Verify(token string) (*User, error)
}

type service struct {
	tokenGen TokenGenerator
}

func NewService(tg TokenGenerator) Service {
	return &service{tg}
}

func (s *service) Login(username, pass string) (*AuthResponse, error) {
	if username != "admin" || pass != "123" {
		return nil, ErrInvalidCredentials
	}

	user := &User{ID: "aK0o3", Name: "John Doe"}

	token, err := s.getToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.getRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{user, token, refreshToken}, nil
}

func (s *service) Verify(token string) (*User, error) {
	if token != "arandomtokenhere" {
		return nil, ErrInvalidToken
	}
	return &User{ID: "aK0o3", Name: "John Doe"}, nil
func (s *service) getToken(payload Payload) (string, error) {
func (s *service) getToken(user *User) (string, error) {
	payload := s.getPayload(user, time.Now().Add(time.Hour))
	return s.tokenGen.Sign(payload, os.Getenv(JWT_SIGN_SECRET_ENV))
}

func (s *service) getRefreshToken(user *User) (string, error) {
	payload := s.getPayload(user, time.Now().AddDate(1, 0, 0))
	return s.tokenGen.Sign(payload, os.Getenv(JWT_REFRESH_SECRET_ENV))
}

func (s *service) getPayload(user *User, exp time.Time) Payload {
	return Payload{
		"iss":  "auth",
		"exp":  exp,
		"sub":  "admin",
		"aud":  "renting",
		"user": user,
	}
}
