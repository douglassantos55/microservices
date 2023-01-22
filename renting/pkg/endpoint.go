package pkg

import (
	"context"
	"math"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	Create endpoint.Endpoint
	List   endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	return Set{
		Create: createRentEndpoint(svc),
		List:   createListEndpoint(svc),
	}
}

func createRentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.CreateRent(r.(Rent))
	}
}

func createListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		pagination := r.(Pagination)
		rents, total, err := svc.ListRents(pagination.Page, pagination.PerPage)
		if err != nil {
			return nil, err
		}

		pages := int64(math.Max(1, math.Round(float64(total/pagination.PerPage))))

		items := make([]any, len(rents))
		for i, rent := range rents {
			items[i] = rent
		}

		return ListResult{items, pages, total}, nil
	}
}

type Pagination struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

type ListResult struct {
	Items      []any `json:"items"`
	TotalPages int64 `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}
