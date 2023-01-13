package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreatePaymentMethod endpoint.Endpoint
	ListPaymentMethods  endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	return Set{
		CreatePaymentMethod: makeCreatePaymentMethodEndpoint(svc),
		ListPaymentMethods:  makeListPaymentMethodsEndpoint(svc),
	}
}

func makeCreatePaymentMethodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		method := r.(PaymentMethod)
		return svc.CreatePaymentMethod(method)
	}
}

func makeListPaymentMethodsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.ListPaymentMethods()
	}
}
