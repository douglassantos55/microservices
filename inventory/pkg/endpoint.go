package pkg

import (
	"context"
	"math"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	Get          endpoint.Endpoint
	Create       endpoint.Endpoint
	List         endpoint.Endpoint
	Update       endpoint.Endpoint
	Delete       endpoint.Endpoint
	ReduceStock  endpoint.Endpoint
	RestoreStock endpoint.Endpoint
}

func NewSet(svc Service) Set {
	return Set{
		Get:          makeGetEndpoint(svc),
		Create:       makeCreateEndpoint(svc),
		List:         makeListEndpoint(svc),
		Update:       makeUpdateEndpoint(svc),
		Delete:       makeDeleteEndpoint(svc),
		ReduceStock:  makeReduceStockEndpoint(svc),
		RestoreStock: makeRestoreStockEndpoint(svc),
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

func makeDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return nil, svc.DeleteEquipment(r.(string))
	}
}

func makeGetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetEquipment(r.(string))
	}
}

func makeReduceStockEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(ReduceStockRequest)
		return nil, svc.ReduceStock(req.Equip, req.Qty)
	}
}

func makeRestoreStockEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(RestoreStockRequest)
		return nil, svc.RestoreStock(req.Equip, req.Qty)
	}
}
