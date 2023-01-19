package pkg

import (
	"context"

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
