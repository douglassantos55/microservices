package pkg

import (
	"context"

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
