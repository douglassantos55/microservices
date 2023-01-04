package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPHandler(svc Service) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodGet, "/verify", httptransport.NewServer(
		makeVerifyEndpoint(svc),
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
	))

	return router
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (any, error) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return nil, err
	}
	return credentials, nil
}
