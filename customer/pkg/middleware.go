package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func verifyMiddleware(service string) endpoint.Middleware {
	verify := makeVerifyEndpoint(service)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			_, err := verify(ctx, r)
			if err != nil {
				return nil, err
			}
			return next(ctx, r)
		}
	}
}
