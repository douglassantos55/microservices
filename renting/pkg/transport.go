package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPServer(endpoints Set) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		endpoints.Create,
		decodeCreateRentRequest,
		httptransport.EncodeJSONResponse,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		endpoints.List,
		decodeListRequest,
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

func decodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	params := r.URL.Query()
	page, err := strconv.ParseInt(params.Get("page"), 0, 0)
	if err != nil {
		page = 1
	}

	perPage, err := strconv.ParseInt(params.Get("per_page"), 0, 0)
	if err != nil {
		perPage = 50
	}

	return Pagination{page - 1, perPage}, nil
}
