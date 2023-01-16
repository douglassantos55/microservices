package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPServer(endpoints Set) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		endpoints.CreateRent,
		decodeCreateRentRequest,
		httptransport.EncodeJSONResponse,
	))

	return router
}

func decodeCreateRentRequest(ctx context.Context, r *http.Request) (any, error) {
	var rent Rent
	if err := json.NewDecoder(r.Body).Decode(&rent); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return rent, nil
}
