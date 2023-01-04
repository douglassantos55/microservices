package pkg

import (
	"context"
	"errors"
	"time"

	"api.example.com/microservices.git/auth/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func makeVerifyEndpoint(authService string) endpoint.Endpoint {
	conn, err := grpc.Dial(
		authService,
		grpc.WithInsecure(),
		grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		panic(err)
	}

	return grpctransport.NewClient(
		conn,
		"pb.Auth",
		"Verify",
		encodeVerifyRequest,
		decodeVerifyResponse,
		&proto.VerifyReply{},
	).Endpoint()
}

func encodeVerifyRequest(ctx context.Context, r any) (any, error) {
	token := ctx.Value(jwt.JWTContextKey).(string)
	return &proto.Token{Token: token}, nil
}

func decodeVerifyResponse(ctx context.Context, r any) (any, error) {
	res := r.(*proto.VerifyReply)
	if res.GetErr() != "" {
		return nil, errors.New(res.GetErr())
	}
	return VerifyResponse{res.User.String(), res.Err}, nil
}

type VerifyResponse struct {
	User string `json:"user"`
	Err  string `json:"err,omitempty"`
}
