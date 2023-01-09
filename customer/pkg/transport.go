package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPHandler(set Set) http.Handler {
	router := httprouter.New()

	opts := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	router.Handler(http.MethodGet, "/:id", httptransport.NewServer(
		set.Get,
		GetURLParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		opts...,
	))

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		set.Create,
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
		opts...,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		set.List,
		decodeListRequest,
		httptransport.EncodeJSONResponse,
		opts...,
	))

	router.Handler(http.MethodPut, "/:id", httptransport.NewServer(
		set.Update,
		decodeUpdateRequest,
		httptransport.EncodeJSONResponse,
		opts...,
	))

	router.Handler(http.MethodDelete, "/:id", httptransport.NewServer(
		set.Delete,
		GetURLParamDecoder("id"),
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

func decodeUpdateRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Customer
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"the provided input is invalid, please verify and try again",
		)
	}

	params := httprouter.ParamsFromContext(r.Context())
	return UpdateRequest{ID: params.ByName("id"), Data: data}, nil
}

type UpdateRequest struct {
	ID   string   `json:"id"`
	Data Customer `json:"data"`
}

func GetURLParamDecoder(param string) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (any, error) {
		params := httprouter.ParamsFromContext(r.Context())
		return params.ByName(param), nil
	}
}
