package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreatePaymentMethod endpoint.Endpoint
	ListPaymentMethods  endpoint.Endpoint
	UpdatePaymentMethod endpoint.Endpoint
	DeletePaymentMethod endpoint.Endpoint
	GetPaymentMethod    endpoint.Endpoint

	CreatePaymentType endpoint.Endpoint
	ListPaymentTypes  endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	return Set{
		CreatePaymentMethod: makeCreatePaymentMethodEndpoint(svc),
		ListPaymentMethods:  makeListPaymentMethodsEndpoint(svc),
		UpdatePaymentMethod: makeUpdatePaymentMethodEndpoint(svc),
		DeletePaymentMethod: makeDeletePaymentMethodEndpoint(svc),
		GetPaymentMethod:    makeGetPaymentMethodEndpoint(svc),

		CreatePaymentType: makeCreatePaymentTypeEndpoint(svc),
		ListPaymentTypes:  makeListPaymentTypesEndpoint(svc),
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

func makeUpdatePaymentMethodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdatePaymentMethodRequest)
		return svc.UpdatePaymentMethod(req.ID, req.Data)
	}
}

func makeDeletePaymentMethodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return nil, svc.DeletePaymentMethod(r.(string))
	}
}

func makeGetPaymentMethodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetPaymentMethod(r.(string))
	}
}

func makeCreatePaymentTypeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		data := r.(PaymentType)
		return svc.CreatePaymentType(data)
	}
}

func makeListPaymentTypesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.ListPaymentTypes()
	}
}
