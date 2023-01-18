package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/renting/proto"
)

func WithPaymentTypeEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withPaymentType := withPaymentTypeMiddleware(cc)
	return Set{
		CreateRent: withPaymentType(endpoints.CreateRent),
	}
}

func withPaymentTypeMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getPaymentType := getPaymentTypeEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				paymentType, err := getPaymentType(ctx, rent.PaymentTypeID)
				if err == nil {
					rent.PaymentType = paymentType.(*PaymentType)
				}
				return rent, nil
			}

			return res, err
		}
	}
}

func getPaymentTypeEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Payment",
		"GetType",
		encodeRequest,
		decodePaymentType,
		&proto.TypeReply{},
	).Endpoint()
}

func encodeRequest(ctx context.Context, r any) (any, error) {
	return &proto.GetRequest{Id: r.(string)}, nil
}

func decodePaymentType(ctx context.Context, r any) (any, error) {
	res := r.(*proto.TypeReply)
	paymentType := res.GetType()

	return &PaymentType{
		ID:   paymentType.GetId(),
		Name: paymentType.GetName(),
	}, nil
}

func WithPaymentMethodEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withPaymentMethod := withPaymentMethodMiddleware(cc)

	return Set{
		CreateRent: withPaymentMethod(endpoints.CreateRent),
	}
}

func withPaymentMethodMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getPaymentMethod := getPaymentMethodEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				method, err := getPaymentMethod(ctx, rent.PaymentMethodID)
				if err == nil {
					rent.PaymentMethod = method.(*PaymentMethod)
				}
				return rent, nil
			}

			return res, nil
		}
	}
}

func getPaymentMethodEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Payment",
		"GetMethod",
		encodeRequest,
		decodePaymentMethod,
		&proto.MethodReply{},
	).Endpoint()
}

func decodePaymentMethod(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.MethodReply)

	return &PaymentMethod{
		ID:   reply.Method.GetId(),
		Name: reply.Method.GetName(),
	}, nil
}

func WithPaymentConditionEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withPaymentCondition := withPaymentConditionMiddleware(cc)

	return Set{
		CreateRent: withPaymentCondition(endpoints.CreateRent),
	}
}

func withPaymentConditionMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getPaymentCondition := getPaymentConditionEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				condition, err := getPaymentCondition(ctx, rent.PaymentConditionID)
				if err == nil {
					rent.PaymentCondition = condition.(*PaymentCondition)
				}
				return rent, nil
			}

			return res, nil
		}
	}
}

func getPaymentConditionEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Payment",
		"GetCondition",
		encodeRequest,
		decodePaymentCondition,
		&proto.ConditionReply{},
	).Endpoint()
}

func decodePaymentCondition(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.ConditionReply)

	condition := reply.GetCondition()
	paymentType := condition.GetPaymentType()

	return &PaymentCondition{
		ID:        condition.GetId(),
		Name:      condition.GetName(),
		Increment: condition.GetIncrement(),
		PaymentType: &PaymentType{
			ID:   paymentType.GetId(),
			Name: paymentType.GetName(),
		},
		Installments: condition.GetInstallments(),
	}, nil
}

func WithCustomerEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withCustomer := withCustomerMiddleware(cc)

	return Set{
		CreateRent: withCustomer(endpoints.CreateRent),
	}
}

func withCustomerMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getCustomer := getCustomerEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				customer, err := getCustomer(ctx, rent.CustomerID)
				if err == nil {
					rent.Customer = customer.(*Customer)
				}
				return rent, nil
			}

			return res, nil
		}
	}
}

func getCustomerEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Customer",
		"Get",
		encodeRequest,
		decodeCustomer,
		&proto.Customer{},
	).Endpoint()
}

func decodeCustomer(ctx context.Context, r any) (any, error) {
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
