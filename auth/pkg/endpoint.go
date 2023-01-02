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
		cred := r.(Credentials)
		response, err := svc.Login(cred.User, cred.Pass)

		var error string
		if err != nil {
			error = err.Error()
		}

		return LoginResponse{
			User:    response.User,
			Token:   response.Token,
			Refresh: response.Refresh,
			Err:     error,
		}, nil
	}
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginResponse struct {
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	Err     string `json:"err,omitempty"`
}

func (r LoginResponse) StatusCode() int {
	if r.Err != "" {
		return http.StatusBadRequest
	}
	return http.StatusOK
}
