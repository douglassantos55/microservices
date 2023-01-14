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

	makePaymentMethodRoutes("/methods", router, endpoints)
	makePaymentTypeRoutes("/types", router, endpoints)

	return router
}

func makePaymentMethodRoutes(prefix string, router *httprouter.Router, endpoints Set) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreatePaymentMethod,
		decodeCreatePaymentMethodRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodGet, prefix+"/", httptransport.NewServer(
		endpoints.ListPaymentMethods,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodPut, prefix+"/:id", httptransport.NewServer(
		endpoints.UpdatePaymentMethod,
		decodeUpdatePaymentMethodRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodDelete, prefix+"/:id", httptransport.NewServer(
		endpoints.DeletePaymentMethod,
		GetRouteParamDecoder("id"),
		encodeDeleteResponse,
	))

	router.Handler(http.MethodGet, prefix+"/:id", httptransport.NewServer(
		endpoints.GetPaymentMethod,
		GetRouteParamDecoder("id"),
		httptransport.EncodeJSONResponse,
	))
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

func GetRouteParamDecoder(param string) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (any, error) {
		params := httprouter.ParamsFromContext(r.Context())
		return params.ByName(param), nil
	}
}

func encodeDeleteResponse(ctx context.Context, r http.ResponseWriter, res any) error {
	r.WriteHeader(http.StatusNoContent)
	return nil
}

func makePaymentTypeRoutes(prefix string, router *httprouter.Router, endpoints Set) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreatePaymentType,
		decodeCreatePaymentTypeRequest,
		httptransport.EncodeJSONResponse,
	))
}

func decodeCreatePaymentTypeRequest(ctx context.Context, r *http.Request) (any, error) {
	var paymentType PaymentType
	if err := json.NewDecoder(r.Body).Decode(&paymentType); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return paymentType, nil
}
