package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeVerifyEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, token any) (any, error) {
		user, err := svc.Verify(token.(string))
		return VerifyResponse{user, err}, nil
	}
}

type VerifyResponse struct {
	User *User `json:"user"`
	Err  error `json:"-"`
}

func (r VerifyResponse) Failed() error {
	return r.Err
}

func makeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		credentials := r.(Credentials)
		token, err := svc.Login(credentials.User, credentials.Pass)

		return LoginResponse{token, err}, nil
	}
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Err   error  `json:"-"`
}

func (r LoginResponse) Failed() error {
	return r.Err
}
