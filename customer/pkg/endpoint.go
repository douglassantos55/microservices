package pkg

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	Get    endpoint.Endpoint
	List   endpoint.Endpoint
	Create endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

func NewSet(svc Service) Set {
	return Set{
		Get:    makeGetEndpoint(svc),
		List:   makeListEndpoint(svc),
		Create: makeCreateEndpoint(svc),
		Update: makeUpdateEndpoint(svc),
		Delete: makeDeleteEndpoint(svc),
	}
}

func makeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.Create(r.(Customer))
	}
}

func makeGetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.Get(r.(string))
	}
}

func makeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		pagination := r.(Pagination)
		return svc.List(pagination.Page, pagination.PerPage)
	}
}

func makeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdateRequest)
		return svc.Update(req.ID, req.Data)
	}
}

func makeDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		if err := svc.Delete(r.(string)); err != nil {
			return nil, NewError(
				http.StatusInternalServerError,
				"could not delete customer",
				"something went wrong while deleting customer",
			)
		}
		return DeleteResponse{}, nil
	}
}

type DeleteResponse struct{}

func (r DeleteResponse) StatusCode() int {
	return http.StatusNoContent
}
