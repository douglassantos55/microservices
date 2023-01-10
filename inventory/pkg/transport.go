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
		endpoints.Create,
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
	))

	return router
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	var equipment Equipment
	if err := json.NewDecoder(r.Body).Decode(&equipment); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data and try again",
		)
	}
	return equipment, nil
}
