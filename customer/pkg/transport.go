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

func MakeHTTPHandler(svc Service, cc *grpc.ClientConn) http.Handler {
	router := httprouter.New()
	verify := verifyMiddleware(cc)

	opts := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		verify(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
		opts...,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		verify(makeListEndpoint(svc)),
		decodeListRequest,
		httptransport.EncodeJSONResponse,
		opts...,
	))

	return router
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

func decodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	params := r.URL.Query()
	page, err := strconv.ParseInt(params.Get("page"), 0, 0)
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err := strconv.ParseInt(params.Get("per_page"), 0, 0)
	if err != nil {
		perPage = 50
	}

	return Pagination{Page: page - 1, PerPage: perPage}, nil
}

type Pagination struct {
	Page    int64
	PerPage int64
}
