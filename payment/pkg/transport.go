package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPHandler(endpoints Set) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		endpoints.CreatePaymentMethod,
		decodeCreatePaymentMethodRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		endpoints.ListPaymentMethods,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
	))

	return router
}

func decodeCreatePaymentMethodRequest(ctx context.Context, r *http.Request) (any, error) {
	var method PaymentMethod

	if err := json.NewDecoder(r.Body).Decode(&method); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return method, nil
}
