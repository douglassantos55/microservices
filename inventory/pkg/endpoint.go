package pkg

import (
	"context"
	"math"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	Create endpoint.Endpoint
	List   endpoint.Endpoint
	Update endpoint.Endpoint
}

func NewSet(svc Service) Set {
	return Set{
		Create: makeCreateEndpoint(svc),
		List:   makeListEndpoint(svc),
		Update: makeUpdateEndpoint(svc),
	}
}

func makeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		equipment := r.(Equipment)
		return svc.CreateEquipment(equipment)
	}
}

func makeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		pagination := r.(Pagination)
		equipment, total, err := svc.ListEquipment(pagination.Page, pagination.PerPage)
		if err != nil {
			return nil, err
		}

		items := make([]any, len(equipment))
		for i, equip := range equipment {
			items[i] = equip
		}

		totalPages := int(math.Max(1, math.Round(float64(total)/float64(pagination.PerPage))))

		return ListResult{
			Items:      items,
			TotalItems: total,
			TotalPages: totalPages,
		}, nil
	}
}

type ListResult struct {
	Items      []any `json:"items"`
	TotalItems int   `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func makeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdateRequest)
		return svc.UpdateEquipment(req.ID, req.Data)
	}
}
