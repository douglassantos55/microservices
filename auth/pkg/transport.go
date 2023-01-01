package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	"api.example.com/auth/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(svc Service) http.Handler {
	return httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		httptransport.EncodeJSONResponse,
	)
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (any, error) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return nil, err
	}
	return credentials, nil
}

func NewGRPCHandler(svc Service) grpctransport.Handler {
	return grpctransport.NewServer(
		makeVerifyEndpoint(svc),
		decodeVerifyRequest,
		encodeVerifyResponse,
	)
}

func decodeVerifyRequest(ctx context.Context, r any) (any, error) {
	token := r.(*proto.Token)
	return token.GetToken(), nil
}

func encodeVerifyResponse(ctx context.Context, r any) (any, error) {
	var user *proto.User
	reply := r.(VerifyResponse)

	if reply.User != nil {
		user.Id = reply.User.ID
		user.Name = reply.User.Name
	}

	return proto.VerifyReply{User: user, Err: reply.Err.Error()}, nil
}
