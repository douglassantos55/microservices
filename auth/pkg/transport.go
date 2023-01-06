package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"reconcip.com.br/microservices/auth/proto"
)

func NewHTTPHandler(svc Service) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		httptransport.EncodeJSONResponse,
	))

	return router
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (any, error) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"the provided input is invalid, please verify and try again",
		)
	}
	return credentials, nil
}

type grpcServer struct {
	proto.UnimplementedAuthServer
	verify grpc.Handler
}

func NewGRPCServer(svc Service) proto.AuthServer {
	return &grpcServer{
		verify: grpc.NewServer(
			makeVerifyEndpoint(svc),
			nopGRPCRequestDecoder,
			encodeVerifyResponse,
			grpc.ServerBefore(jwt.GRPCToContext()),
		),
	}
}

func (s *grpcServer) Verify(ctx context.Context, r *proto.Token) (*proto.VerifyReply, error) {
	_, rep, err := s.verify.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.VerifyReply), nil
}

// NopGRCPRequestDecoder is a DecodeRequestFunc that can be used for requests
// that do not need to be decoded, and simply returns nil, nil.
func nopGRPCRequestDecoder(ctx context.Context, r any) (any, error) {
	return nil, nil
}

func encodeVerifyResponse(ctx context.Context, r any) (any, error) {
	reply := r.(*User)
	user := &proto.User{Id: reply.ID, Name: reply.Name}
	return &proto.VerifyReply{User: user}, nil
}
