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
	getPaymentMethod := paymentMethodEndpoint(cc)

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

func paymentMethodEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
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
