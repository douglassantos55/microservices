package pkg

import (
	"context"
	"math"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		data := r.(Supplier)
		supplier, err := svc.Create(data)
		if err != nil {
			return nil, err
		}
		return supplier, nil
	}
}

func makeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		pagination := r.(Pagination)
		suppliers, totalItems, err := svc.List(pagination.Page, pagination.PerPage)

		if err != nil {
			return nil, NewError(
				http.StatusInternalServerError,
				"error fetching suppliers",
				"something went wrong while fetching suppliers",
			)
		}

		totalPages := int64(math.Max(1, math.Round(float64(totalItems)/float64(pagination.PerPage))))
		if pagination.Page >= int64(totalPages) {
			return nil, NewError(
				http.StatusBadRequest,
				"invalid page",
				"page exceeds the maximum number of pages available",
			)
		}

		items := make([]any, len(suppliers))
		for i, supplier := range suppliers {
			items[i] = supplier
		}

		return ListResult{
			Items:      items,
			TotalItems: totalItems,
			TotalPages: totalPages,
		}, nil
	}
}

type ListResult struct {
	Items      []any `json:"items"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

func makeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdateRequest)
		return svc.Update(req.ID, req.Data)
	}
}
