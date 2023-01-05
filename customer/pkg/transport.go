package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(svc Service, authService string) http.Handler {
	verify := verifyMiddleware(authService)

	return httptransport.NewServer(
		verify(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Customer
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"the provided input is invalid, please verify and try again",
		)
	}
	return data, nil
}
