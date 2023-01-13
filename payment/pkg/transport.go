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

	router.Handler(http.MethodPut, "/:id", httptransport.NewServer(
		endpoints.UpdatePaymentMethod,
		decodeUpdatePaymentMethodRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodDelete, "/:id", httptransport.NewServer(
		endpoints.DeletePaymentMethod,
		decodeDeletePaymentMethodRequest,
		encodeDeleteResponse,
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

func decodeUpdatePaymentMethodRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var method PaymentMethod
	if err := json.NewDecoder(r.Body).Decode(&method); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return UpdatePaymentMethodRequest{
		ID:   params.ByName("id"),
		Data: method,
	}, nil
}

type UpdatePaymentMethodRequest struct {
	ID   string
	Data PaymentMethod
}

func decodeDeletePaymentMethodRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName("id"), nil
}

func encodeDeleteResponse(ctx context.Context, r http.ResponseWriter, res any) error {
	r.WriteHeader(http.StatusNoContent)
    return nil
}
