package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/renting/proto"
)

func WithPaymentTypeEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withPaymentType := WithPaymentTypeMiddleware(cc)
	return Set{
		CreateRent: withPaymentType(endpoints.CreateRent),
	}
}

func WithPaymentTypeMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
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
		decodeResponse,
		&proto.TypeReply{},
	).Endpoint()
}

func encodeRequest(ctx context.Context, r any) (any, error) {
	return &proto.GetRequest{Id: r.(string)}, nil
}

func decodeResponse(ctx context.Context, r any) (any, error) {
	res := r.(*proto.TypeReply)
	paymentType := res.GetType()

	return &PaymentType{
		ID:   paymentType.GetId(),
		Name: paymentType.GetName(),
	}, nil
}
