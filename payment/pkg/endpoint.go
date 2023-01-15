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
	UpdatePaymentType endpoint.Endpoint
	DeletePaymentType endpoint.Endpoint
	GetPaymentType    endpoint.Endpoint

	CreatePaymentCondition endpoint.Endpoint
	ListPaymentConditions  endpoint.Endpoint
	UpdatePaymentCondition endpoint.Endpoint
	DeletePaymentCondition endpoint.Endpoint
	GetPaymentCondition    endpoint.Endpoint
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
		UpdatePaymentType: makeUpdatePaymentTypeEndpoint(svc),
		DeletePaymentType: makeDeletePaymentTypeEndpoint(svc),
		GetPaymentType:    makeGetPaymentTypeEndpoint(svc),

		CreatePaymentCondition: makeCreatePaymentConditionEndpoint(svc),
		ListPaymentConditions:  makeListPaymentConditionsEndpoint(svc),
		UpdatePaymentCondition: makeUpdatePaymentConditionEndpoint(svc),
		DeletePaymentCondition: makeDeletePaymentConditionEndpoint(svc),
		GetPaymentCondition:    makeGetPaymentConditionEndpoint(svc),
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

func makeUpdatePaymentTypeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdatePaymentTypeRequest)
		return svc.UpdatePaymentType(req.ID, req.Data)
	}
}

func makeDeletePaymentTypeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return nil, svc.DeletePaymentType(r.(string))
	}
}

func makeGetPaymentTypeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetPaymentType(r.(string))
	}
}

func makeCreatePaymentConditionEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		condition := r.(Condition)
		return svc.CreatePaymentCondition(condition)
	}
}

func makeListPaymentConditionsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.ListPaymentConditions()
	}
}

func makeUpdatePaymentConditionEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdateConditionRequest)
		return svc.UpdatePaymentCondition(req.ID, req.Data)
	}
}

func makeDeletePaymentConditionEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return nil, svc.DeletePaymentCondition(r.(string))
	}
}

func makeGetPaymentConditionEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetPaymentCondition(r.(string))
	}
}
