package pkg

import (
	"context"
	"log"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/payment/proto"
)

func VerifyEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	verify := verifyMiddleware(cc)

	return Set{
		CreatePaymentMethod: verify(endpoints.CreatePaymentMethod),
		ListPaymentMethods:  verify(endpoints.ListPaymentMethods),
		UpdatePaymentMethod: verify(endpoints.UpdatePaymentMethod),
		DeletePaymentMethod: verify(endpoints.DeletePaymentMethod),
		GetPaymentMethod:    verify(endpoints.GetPaymentMethod),

		CreatePaymentType: verify(endpoints.CreatePaymentType),
		ListPaymentTypes:  verify(endpoints.ListPaymentTypes),
		UpdatePaymentType: verify(endpoints.UpdatePaymentType),
		DeletePaymentType: verify(endpoints.DeletePaymentType),
		GetPaymentType:    verify(endpoints.GetPaymentType),

		CreatePaymentCondition: verify(endpoints.CreatePaymentCondition),
		ListPaymentConditions:  verify(endpoints.ListPaymentConditions),
		UpdatePaymentCondition: verify(endpoints.UpdatePaymentCondition),
		DeletePaymentCondition: verify(endpoints.DeletePaymentCondition),
		GetPaymentCondition:    verify(endpoints.GetPaymentCondition),

		CreateInvoice: verify(endpoints.CreateInvoice),
		ListInvoices:  verify(endpoints.ListInvoices),
		UpdateInvoice: verify(endpoints.UpdateInvoice),
		DeleteInvoice: verify(endpoints.DeleteInvoice),
		GetInvoice:    verify(endpoints.GetInvoice),
	}
}

func verifyMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	verify := verifyEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			if _, err := verify(ctx, r); err != nil {
				return nil, err
			}
			return next(ctx, r)
		}
	}
}

func verifyEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Auth",
		"Verify",
		encodeVerifyRequest,
		decodeVerifyResponse,
		&proto.VerifyReply{},
		grpctransport.ClientBefore(jwt.ContextToGRPC()),
	).Endpoint()
}

func encodeVerifyRequest(ctx context.Context, r any) (any, error) {
	return nil, nil
}

func decodeVerifyResponse(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.VerifyReply)
	if reply.Err != nil {
		return nil, NewErrorFromReply(reply.Err)
	}

	var user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	user.ID = reply.User.Id
	user.Name = reply.User.Name

	return user, nil

}

func CustomerEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withCustomer := withCustomerMiddleware(cc)

	return Set{
		CreatePaymentMethod: endpoints.CreatePaymentMethod,
		ListPaymentMethods:  endpoints.ListPaymentMethods,
		UpdatePaymentMethod: endpoints.UpdatePaymentMethod,
		DeletePaymentMethod: endpoints.DeletePaymentMethod,
		GetPaymentMethod:    endpoints.GetPaymentMethod,

		CreatePaymentType: endpoints.CreatePaymentType,
		ListPaymentTypes:  endpoints.ListPaymentTypes,
		UpdatePaymentType: endpoints.UpdatePaymentType,
		DeletePaymentType: endpoints.DeletePaymentType,
		GetPaymentType:    endpoints.GetPaymentType,

		CreatePaymentCondition: endpoints.CreatePaymentCondition,
		ListPaymentConditions:  endpoints.ListPaymentConditions,
		UpdatePaymentCondition: endpoints.UpdatePaymentCondition,
		DeletePaymentCondition: endpoints.DeletePaymentCondition,
		GetPaymentCondition:    endpoints.GetPaymentCondition,

		CreateInvoice: withCustomer(endpoints.CreateInvoice),
		ListInvoices:  withCustomer(endpoints.ListInvoices),
		UpdateInvoice: withCustomer(endpoints.UpdateInvoice),
		DeleteInvoice: withCustomer(endpoints.DeleteInvoice),
		GetInvoice:    withCustomer(endpoints.GetInvoice),
	}
}

func withCustomerMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getCustomer := getCustomerProxy(cc)

	appendCustomer := func(ctx context.Context, invoice *Invoice) {
		customer, err := getCustomer(ctx, invoice.CustomerID)
		if err != nil {
			log.Printf("could not find customer %v: %v", invoice.CustomerID, err)
		} else {
			invoice.Customer = customer.(*Customer)
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if invoice, ok := res.(*Invoice); ok {
				appendCustomer(ctx, invoice)
			}

			if result, ok := res.(ListResult); ok {
				for _, item := range result.Items {
					appendCustomer(ctx, item.(*Invoice))
				}
			}

			return res, nil
		}
	}
}

func getCustomerProxy(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Customer",
		"Get",
		encodeGetRequest,
		decodeCustomerResponse,
		&proto.Customer{},
	).Endpoint()
}

func encodeGetRequest(ctx context.Context, r any) (any, error) {
	return &proto.GetRequest{Id: r.(string)}, nil
}

func decodeCustomerResponse(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.Customer)

	return &Customer{
		ID:        reply.GetId(),
		Name:      reply.GetName(),
		Email:     reply.GetEmail(),
		CpfCnpj:   reply.GetCpfCnpj(),
		RgInscEst: reply.GetRgInscEst(),
		Phone:     reply.GetPhone(),
		Cellphone: reply.GetCellphone(),
	}, nil
}
