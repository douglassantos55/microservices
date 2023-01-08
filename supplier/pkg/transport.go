package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		makeListEndpoint(svc),
		decodeListRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodPut, "/:id", httptransport.NewServer(
		makeUpdateEndpoint(svc),
		decodeUpdateRequest,
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
