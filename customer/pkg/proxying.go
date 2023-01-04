package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeVerifyEndpoint(authService string) endpoint.Endpoint {
	verifyUrl, err := url.Parse("http://" + authService + "/verify")
	if err != nil {
		panic(err)
	}

	return httptransport.NewClient(
		"GET",
		verifyUrl,
		encodeVerifyRequest,
		decodeVerifyResponse,
	).Endpoint()
}

func encodeVerifyRequest(ctx context.Context, req *http.Request, data any) error {
	token, ok := ctx.Value(jwt.JWTContextKey).(string)
	if !ok {
		return jwt.ErrTokenContextMissing
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func decodeVerifyResponse(ctx context.Context, r *http.Response) (any, error) {
	var user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}
