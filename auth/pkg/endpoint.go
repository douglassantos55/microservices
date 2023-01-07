package pkg

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

func makeVerifyEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		token, ok := ctx.Value(jwt.JWTContextKey).(string)
		if !ok {
			return VerifyResponse{nil, NewError(
				http.StatusBadRequest,
				"Empty authorization token",
				"could not find token in authorization header",
			)}, nil
		}
		user, err := svc.Verify(token)
		return VerifyResponse{user, err}, nil
	}
}

type VerifyResponse struct {
	User *User `json:"user"`
	Err  Error `json:"-"`
}

func (r VerifyResponse) Failed() error {
	return r.Err
}

func makeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		cred := r.(Credentials)
		response, err := svc.Login(cred.User, cred.Pass)

		if err != nil {
			return nil, err
		}

		return LoginResponse{
			User:    response.User,
			Token:   response.Token,
			Refresh: response.Refresh,
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
}
