package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	Create endpoint.Endpoint
}

func NewSet(svc Service) Set {
	return Set{
		Create: makeCreateEndpoint(svc),
	}
}

func makeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		equipment := r.(Equipment)
		return svc.CreateEquipment(equipment)
	}
}
