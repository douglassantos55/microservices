package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GetQuotes endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	return Set{
		GetQuotes: makeGetQuotesEndpoint(svc),
	}
}

func makeGetQuotesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(GetQuotesRequest)
		return svc.GetQuotes(req.Origin, req.Dest, req.Items)
	}
}
