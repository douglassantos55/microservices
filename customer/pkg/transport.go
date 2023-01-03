package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(svc Service) http.Handler {
	return httptransport.NewServer(
		makeCreateEndpoint(svc),
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
	)
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Customer
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, ErrBadRequest{}
	}
	return data, nil
}

type ErrBadRequest struct{}

func (e ErrBadRequest) Error() string {
	return "invalid request"
}

func (e ErrBadRequest) StatusCode() int {
	return http.StatusBadRequest
}
