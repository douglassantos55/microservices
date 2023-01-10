package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"reconcip.com.br/microservices/supplier/proto"
)

func NewHTTPServer(set Set) http.Handler {
	router := httprouter.New()

	options := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		set.Create,
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		set.List,
		decodeListRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodPut, "/:id", httptransport.NewServer(
		set.Update,
		decodeUpdateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodDelete, "/:id", httptransport.NewServer(
		set.Delete,
		GetUrlParamDecoder("id"),
		encodeDeleteResponse,
		options...,
	))

	router.Handler(http.MethodGet, "/:id", httptransport.NewServer(
		set.Get,
		GetUrlParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		options...,
	))

	return router
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	var supplier Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data",
		)
	}
	return supplier, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	params := r.URL.Query()
	page, err := strconv.ParseInt(params.Get("page"), 0, 0)
	if err != nil || page <= 0 {
		page = 1
	}
	perPage, err := strconv.ParseInt(params.Get("per_page"), 0, 0)
	if err != nil || perPage <= 0 {
		perPage = 50
	}
	return Pagination{page - 1, perPage}, nil
}

type Pagination struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Supplier
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data and try again",
		)
	}

	params := httprouter.ParamsFromContext(r.Context())
	return UpdateRequest{params.ByName("id"), data}, nil
}

type UpdateRequest struct {
	ID   string   `json:"id"`
	Data Supplier `json:"data"`
}

func GetUrlParamDecoder(param string) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (any, error) {
		params := httprouter.ParamsFromContext(r.Context())
		return params.ByName(param), nil
	}
}

func encodeDeleteResponse(ctx context.Context, res http.ResponseWriter, r any) error {
	res.WriteHeader(http.StatusNoContent)
	return nil
}

type grpcServer struct {
	proto.UnimplementedSupplierServiceServer
	get grpc.Handler
}

func NewGRPCServer(endpoints Set) proto.SupplierServiceServer {
	return &grpcServer{
		get: grpc.NewServer(
			endpoints.Get,
			decodeGRPCGetRequest,
			encodeGRPCGetResponse,
		),
	}
}

func (s *grpcServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.Supplier, error) {
	_, reply, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return reply.(*proto.Supplier), nil
}

func decodeGRPCGetRequest(ctx context.Context, req any) (any, error) {
	request := req.(*proto.GetRequest)
	return request.SupplierID, nil
}

func encodeGRPCGetResponse(ctx context.Context, res any) (any, error) {
	supplier := res.(*Supplier)

	return &proto.Supplier{
		Id:         supplier.ID,
		SocialName: supplier.SocialName,
		LegalName:  supplier.LegalName,
		Email:      supplier.Email,
		Website:    supplier.Website,
		Cnpj:       supplier.Cnpj,
		InscEst:    supplier.InscEst,
		Phone:      supplier.Phone,
		Address: &proto.Address{
			Street:       supplier.Address.Street,
			Number:       supplier.Address.Number,
			Complement:   supplier.Address.Complement,
			Neighborhood: supplier.Address.Neighborhood,
			City:         supplier.Address.City,
			State:        supplier.Address.State,
			Postcode:     supplier.Address.Postcode,
		},
	}, nil
}
