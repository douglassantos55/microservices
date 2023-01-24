package pkg

import (
	"context"
	"math"

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

	CreateInvoice endpoint.Endpoint
	ListInvoices  endpoint.Endpoint
	UpdateInvoice endpoint.Endpoint
	DeleteInvoice endpoint.Endpoint
	GetInvoice    endpoint.Endpoint
}

func CreateEndpoints(svc Service) Set {
	getType := getTypeMiddleware(svc)
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

		CreatePaymentCondition: getType(makeCreatePaymentConditionEndpoint(svc)),
		ListPaymentConditions:  getType(makeListPaymentConditionsEndpoint(svc)),
		UpdatePaymentCondition: getType(makeUpdatePaymentConditionEndpoint(svc)),
		DeletePaymentCondition: getType(makeDeletePaymentConditionEndpoint(svc)),
		GetPaymentCondition:    getType(makeGetPaymentConditionEndpoint(svc)),

		CreateInvoice: makeCreateInvoiceEndpoint(svc),
		ListInvoices:  makeListInvoicesEndpoint(svc),
		UpdateInvoice: makeUpdateInvoiceEndpoint(svc),
		DeleteInvoice: makeDeleteInvoiceEndpoint(svc),
		GetInvoice:    makeGetInvoiceEndpoint(svc),
	}
}

func makeCreatePaymentMethodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		method := r.(Method)
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
		data := r.(Type)
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

func makeCreateInvoiceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.CreateInvoice(r.(Invoice))
	}
}

func makeListInvoicesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		pagination := r.(Pagination)
		invoices, total, err := svc.ListInvoices(pagination.Page, pagination.PerPage)
		if err != nil {
			return nil, err
		}

		items := make([]any, len(invoices))
		for i, invoice := range invoices {
			items[i] = invoice
		}

		totalPages := int64(math.Max(1, math.Round(float64(total/pagination.PerPage))))

		return ListResult{
			Items:      items,
			TotalItems: total,
			TotalPages: totalPages,
		}, nil
	}
}

type Pagination struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

type ListResult struct {
	Items      []any `json:"items"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

func makeUpdateInvoiceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		req := r.(UpdateInvoiceRequest)
		return svc.UpdateInvoice(req.ID, req.Data)
	}
}

type UpdateInvoiceRequest struct {
	ID   string  `json:"id"`
	Data Invoice `json:"data"`
}

func makeDeleteInvoiceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return nil, svc.DeleteInvoice(r.(string))
	}
}

func makeGetInvoiceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetInvoice(r.(string))
	}
}
