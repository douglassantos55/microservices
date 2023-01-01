package pkg

import "errors"

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
	Login(user, pass string) (string, error)

	// Validates and verifies token
	Verify(token string) (*User, error)
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (s *service) Login(user, pass string) (string, error) {
	if user != "admin" || pass != "123" {
		return "", ErrInvalidCredentials
	}
	return "arandomtokenhere", nil
}

func (s *service) Verify(token string) (*User, error) {
	if token != "arandomtokenhere" {
		return nil, ErrInvalidToken
	}
	return &User{ID: "aK0o3", Name: "John Doe"}, nil
}
