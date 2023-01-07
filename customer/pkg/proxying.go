package pkg

import (
	"context"
	"errors"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/customer/proto"
)

func verifyMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	verify := makeVerifyEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			_, err := verify(ctx, r)
			if err != nil {
				return nil, err
			}
			return next(ctx, r)
		}
	}
}

func makeVerifyEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Auth",
		"Verify",
		nopGRPCRequestEncoder,
		decodeVerifyResponse,
		&proto.VerifyReply{},
		grpctransport.ClientBefore(jwt.ContextToGRPC()),
	).Endpoint()
}

func nopGRPCRequestEncoder(ctx context.Context, r any) (any, error) {
	return nil, nil
}

func decodeVerifyResponse(ctx context.Context, r any) (any, error) {
	rep := r.(*proto.VerifyReply)

	if rep.GetErr() != "" {
		return nil, errors.New(rep.GetErr())
	}

	var user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	user.ID = rep.User.Id
	user.Name = rep.User.Name

	return user, nil
}
