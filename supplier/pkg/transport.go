package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

func NewHTTPServer(svc Service, cc *grpc.ClientConn) http.Handler {
	router := httprouter.New()
	verify := verifyMiddleware(cc)

	options := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		verify(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		verify(makeListEndpoint(svc)),
		decodeListRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodPut, "/:id", httptransport.NewServer(
		verify(makeUpdateEndpoint(svc)),
		decodeUpdateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodDelete, "/:id", httptransport.NewServer(
		verify(makeDeleteEndpoint(svc)),
		GetUrlParamDecoder("id"),
		encodeDeleteResponse,
		options...,
	))

	router.Handler(http.MethodGet, "/:id", httptransport.NewServer(
		verify(makeGetEndpoint(svc)),
		GetUrlParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		options...,
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

func decodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	params := r.URL.Query()
	page, err := strconv.ParseInt(params.Get("page"), 0, 0)
	if err != nil || page <= 0 {
		page = 1
	}
	perPage, err := strconv.ParseInt(params.Get("per_page"), 0, 0)
	if err != nil || perPage <= 0 {
		perPage = 50
	}
	return Pagination{page - 1, perPage}, nil
}

type Pagination struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Supplier
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data and try again",
		)
	}

	params := httprouter.ParamsFromContext(r.Context())
	return UpdateRequest{params.ByName("id"), data}, nil
}

type UpdateRequest struct {
	ID   string   `json:"id"`
	Data Supplier `json:"data"`
}

func GetUrlParamDecoder(param string) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (any, error) {
		params := httprouter.ParamsFromContext(r.Context())
		return params.ByName(param), nil
	}
}

func encodeDeleteResponse(ctx context.Context, res http.ResponseWriter, r any) error {
	res.WriteHeader(http.StatusNoContent)
	return nil
}
