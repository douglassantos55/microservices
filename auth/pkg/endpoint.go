package pkg

import (
	"context"
	"net/http"

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

		var error string
		if err != nil {
			error = err.Error()
		}
		return LoginResponse{token, error}, nil
	}
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"err,omitempty"`
}

func (r LoginResponse) StatusCode() int {
	if r.Err != "" {
		return http.StatusBadRequest
	}
	return http.StatusOK
}
