package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPServer(svc Service) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		makeCreateEndpoint(svc),
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
	))

	return router
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	var supplier Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data",
		)
	}
	return supplier, nil
}
