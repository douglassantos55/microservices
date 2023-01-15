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
	makePaymentConditionRoutes("/conditions", router, endpoints)

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
	var method Method

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

	var method Method
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
	Data Method
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

	router.Handler(http.MethodGet, prefix+"/", httptransport.NewServer(
		endpoints.ListPaymentTypes,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodPut, prefix+"/:id", httptransport.NewServer(
		endpoints.UpdatePaymentType,
		decodeUpdatePaymentTypeRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodDelete, prefix+"/:id", httptransport.NewServer(
		endpoints.DeletePaymentType,
		GetRouteParamDecoder("id"),
		encodeDeleteResponse,
	))

	router.Handler(http.MethodGet, prefix+"/:id", httptransport.NewServer(
		endpoints.GetPaymentType,
		GetRouteParamDecoder("id"),
		httptransport.EncodeJSONResponse,
	))
}

func decodeCreatePaymentTypeRequest(ctx context.Context, r *http.Request) (any, error) {
	var paymentType Type
	if err := json.NewDecoder(r.Body).Decode(&paymentType); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return paymentType, nil
}

func decodeUpdatePaymentTypeRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var method Type
	if err := json.NewDecoder(r.Body).Decode(&method); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return UpdatePaymentTypeRequest{
		ID:   params.ByName("id"),
		Data: method,
	}, nil
}

type UpdatePaymentTypeRequest struct {
	ID   string
	Data Type
}

func makePaymentConditionRoutes(prefix string, router *httprouter.Router, endpoints Set) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreatePaymentCondition,
		decodeCreatePaymentConditionRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodGet, prefix+"/", httptransport.NewServer(
		endpoints.ListPaymentConditions,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodPut, prefix+"/:id", httptransport.NewServer(
		endpoints.UpdatePaymentCondition,
		decodeUpdateConditionRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodDelete, prefix+"/:id", httptransport.NewServer(
		endpoints.DeletePaymentCondition,
		GetRouteParamDecoder("id"),
		encodeDeleteResponse,
	))

	router.Handler(http.MethodGet, prefix+"/:id", httptransport.NewServer(
		endpoints.GetPaymentCondition,
		GetRouteParamDecoder("id"),
		httptransport.EncodeJSONResponse,
	))
}

func decodeCreatePaymentConditionRequest(ctx context.Context, r *http.Request) (any, error) {
	var condition Condition
	if err := json.NewDecoder(r.Body).Decode(&condition); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return condition, nil
}

func decodeUpdateConditionRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var condition Condition
	if err := json.NewDecoder(r.Body).Decode(&condition); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return UpdateConditionRequest{
		ID:   params.ByName("id"),
		Data: condition,
	}, nil
}

type UpdateConditionRequest struct {
	ID   string
	Data Condition
}
