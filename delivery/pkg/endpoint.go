package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GetQuote  endpoint.Endpoint
	GetQuotes endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	return Set{
		GetQuote:  makeGetQuoteEndpoint(svc),
		GetQuotes: makeGetQuotesEndpoint(svc),
	}
}

func makeGetQuoteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(GetQuoteRequest)
		return svc.GetQuote(req.Origin, req.Dest, req.Carrier, req.Items)
	}
}

func makeGetQuotesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(GetQuotesRequest)
		return svc.GetQuotes(req.Origin, req.Dest, req.Items)
	}
}
