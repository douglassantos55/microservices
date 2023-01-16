package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreateRent endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	return Set{
		CreateRent: createRentEndpoint(svc),
	}
}

func createRentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.CreateRent(r.(Rent))
	}
}
