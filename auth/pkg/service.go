package pkg

import (
	"errors"
	"net/http"
	"os"
	"time"
)

const (
	JWT_SIGN_SECRET_ENV    = "JWT_SIGN_SECRET"
	JWT_REFRESH_SECRET_ENV = "JWT_REFRESH_SECRET"
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
	Login(user, pass string) (*AuthResponse, Error)

	// Validates and verifies token
	Verify(token string) (*User, Error)
}

type service struct {
	tokenGen TokenGenerator
}

func NewService(tg TokenGenerator) Service {
	return &service{tg}
}

func (s *service) Login(username, pass string) (*AuthResponse, Error) {
	if username != "admin" || pass != "123" {
		return nil, NewError(
			http.StatusUnprocessableEntity,
			"invalid credentials",
			"could not find user for supplied credentials",
		)
	}

	user := &User{ID: "aK0o3", Name: "John Doe"}

	token, err := s.tokenGen.Sign(
		user,
		time.Now().Add(time.Hour),
		os.Getenv(JWT_SIGN_SECRET_ENV),
	)

	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"could not generate token",
			"something went wrong while generating token, please try again",
		)
	}

	refreshToken, err := s.tokenGen.Sign(
		user,
		time.Now().AddDate(1, 0, 0),
		os.Getenv(JWT_REFRESH_SECRET_ENV),
	)

	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"could not generate token",
			"something went wrong while generating token, please try again",
		)
	}

	return &AuthResponse{user, token, refreshToken}, nil
}

func (s *service) Verify(tokenStr string) (*User, Error) {
	token, err := s.tokenGen.Verify(tokenStr, os.Getenv(JWT_SIGN_SECRET_ENV))
	if err != nil {
		return nil, NewError(
			http.StatusUnauthorized,
			"invalid or expired token",
			"the provided token is either invalid or has expired",
		)
	}

	if !token.IsValid() {
		return nil, NewError(
			http.StatusUnauthorized,
			"invalid token",
			"the provided token is invalid",
		)
	}

	return token.GetUser(), nil
}
