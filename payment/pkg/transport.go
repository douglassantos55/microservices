package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"reconcip.com.br/microservices/payment/proto"
)

type grpcServer struct {
	proto.UnimplementedPaymentServer
	getMethod    grpc.Handler
	getType      grpc.Handler
	getCondition grpc.Handler
}

func NewGRPCServer(endpoints Set) proto.PaymentServer {
	return &grpcServer{
		getMethod: grpc.NewServer(
			endpoints.GetPaymentMethod,
			decodeGRPCGetRequest,
			encodeGRPCMethodReply,
		),
		getType: grpc.NewServer(
			endpoints.GetPaymentType,
			decodeGRPCGetRequest,
			encodeGRPCTypeReply,
		),
		getCondition: grpc.NewServer(
			endpoints.GetPaymentCondition,
			decodeGRPCGetRequest,
			encodeGRPCConditionReply,
		),
	}
}

func (s *grpcServer) GetType(ctx context.Context, r *proto.GetRequest) (*proto.TypeReply, error) {
	_, reply, err := s.getType.ServeGRPC(ctx, r.GetId())
	if err != nil {
		return nil, err
	}
	return reply.(*proto.TypeReply), nil
}

func (s *grpcServer) GetMethod(ctx context.Context, r *proto.GetRequest) (*proto.MethodReply, error) {
	_, reply, err := s.getMethod.ServeGRPC(ctx, r.GetId())
	if err != nil {
		return nil, err
	}
	return reply.(*proto.MethodReply), nil
}

func (s *grpcServer) GetCondition(ctx context.Context, r *proto.GetRequest) (*proto.ConditionReply, error) {
	_, reply, err := s.getCondition.ServeGRPC(ctx, r.GetId())
	if err != nil {
		return nil, err
	}
	return reply.(*proto.ConditionReply), nil
}

func decodeGRPCGetRequest(ctx context.Context, r any) (any, error) {
	return r.(string), nil
}

func encodeGRPCMethodReply(ctx context.Context, r any) (any, error) {
	method := r.(*Method)

	return &proto.MethodReply{
		Method: &proto.Method{
			Id:   method.ID,
			Name: method.Name,
		},
	}, nil
}

func encodeGRPCTypeReply(ctx context.Context, r any) (any, error) {
	paymentType := r.(*Type)

	return &proto.TypeReply{
		Type: &proto.Type{
			Id:   paymentType.ID,
			Name: paymentType.Name,
		},
	}, nil
}

func encodeGRPCConditionReply(ctx context.Context, r any) (any, error) {
	condition := r.(*Condition)

	return &proto.ConditionReply{
		Condition: &proto.Condition{
			Id:           condition.ID,
			Name:         condition.Name,
			Increment:    condition.Increment,
			Installments: condition.Installments,
			PaymentType: &proto.Type{
				Id:   condition.PaymentType.ID,
				Name: condition.PaymentType.Name,
			},
		},
	}, nil
}

func NewHTTPHandler(endpoints Set) http.Handler {
	router := httprouter.New()

	options := httptransport.ServerBefore(
		jwt.HTTPToContext(),
	)

	makePaymentMethodRoutes("/methods", router, endpoints, options)
	makePaymentTypeRoutes("/types", router, endpoints, options)
	makePaymentConditionRoutes("/conditions", router, endpoints, options)
	makeInvoiceRoutes("/invoices", router, endpoints, options)

	return router
}

func makePaymentMethodRoutes(prefix string, router *httprouter.Router, endpoints Set, options httptransport.ServerOption) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreatePaymentMethod,
		decodeCreatePaymentMethodRequest,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodGet, prefix+"/", httptransport.NewServer(
		endpoints.ListPaymentMethods,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodPut, prefix+"/:id", httptransport.NewServer(
		endpoints.UpdatePaymentMethod,
		decodeUpdatePaymentMethodRequest,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodDelete, prefix+"/:id", httptransport.NewServer(
		endpoints.DeletePaymentMethod,
		GetRouteParamDecoder("id"),
		encodeDeleteResponse,
		options,
	))

	router.Handler(http.MethodGet, prefix+"/:id", httptransport.NewServer(
		endpoints.GetPaymentMethod,
		GetRouteParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		options,
	))
}

func decodeCreatePaymentMethodRequest(ctx context.Context, r *http.Request) (any, error) {
	var method Method

	if err := json.NewDecoder(r.Body).Decode(&method); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return method, nil
}

func decodeUpdatePaymentMethodRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var method Method
	if err := json.NewDecoder(r.Body).Decode(&method); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return UpdatePaymentMethodRequest{
		ID:   params.ByName("id"),
		Data: method,
	}, nil
}

type UpdatePaymentMethodRequest struct {
	ID   string
	Data Method
}

func GetRouteParamDecoder(param string) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (any, error) {
		params := httprouter.ParamsFromContext(r.Context())
		return params.ByName(param), nil
	}
}

func encodeDeleteResponse(ctx context.Context, r http.ResponseWriter, res any) error {
	r.WriteHeader(http.StatusNoContent)
	return nil
}

func makePaymentTypeRoutes(prefix string, router *httprouter.Router, endpoints Set, options httptransport.ServerOption) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreatePaymentType,
		decodeCreatePaymentTypeRequest,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodGet, prefix+"/", httptransport.NewServer(
		endpoints.ListPaymentTypes,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodPut, prefix+"/:id", httptransport.NewServer(
		endpoints.UpdatePaymentType,
		decodeUpdatePaymentTypeRequest,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodDelete, prefix+"/:id", httptransport.NewServer(
		endpoints.DeletePaymentType,
		GetRouteParamDecoder("id"),
		encodeDeleteResponse,
		options,
	))

	router.Handler(http.MethodGet, prefix+"/:id", httptransport.NewServer(
		endpoints.GetPaymentType,
		GetRouteParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		options,
	))
}

func decodeCreatePaymentTypeRequest(ctx context.Context, r *http.Request) (any, error) {
	var paymentType Type
	if err := json.NewDecoder(r.Body).Decode(&paymentType); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return paymentType, nil
}

func decodeUpdatePaymentTypeRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var method Type
	if err := json.NewDecoder(r.Body).Decode(&method); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return UpdatePaymentTypeRequest{
		ID:   params.ByName("id"),
		Data: method,
	}, nil
}

type UpdatePaymentTypeRequest struct {
	ID   string
	Data Type
}

func makePaymentConditionRoutes(prefix string, router *httprouter.Router, endpoints Set, options httptransport.ServerOption) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreatePaymentCondition,
		decodeCreatePaymentConditionRequest,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodGet, prefix+"/", httptransport.NewServer(
		endpoints.ListPaymentConditions,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodPut, prefix+"/:id", httptransport.NewServer(
		endpoints.UpdatePaymentCondition,
		decodeUpdateConditionRequest,
		httptransport.EncodeJSONResponse,
		options,
	))

	router.Handler(http.MethodDelete, prefix+"/:id", httptransport.NewServer(
		endpoints.DeletePaymentCondition,
		GetRouteParamDecoder("id"),
		encodeDeleteResponse,
		options,
	))

	router.Handler(http.MethodGet, prefix+"/:id", httptransport.NewServer(
		endpoints.GetPaymentCondition,
		GetRouteParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		options,
	))
}

func decodeCreatePaymentConditionRequest(ctx context.Context, r *http.Request) (any, error) {
	var condition Condition
	if err := json.NewDecoder(r.Body).Decode(&condition); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return condition, nil
}

func decodeUpdateConditionRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var condition Condition
	if err := json.NewDecoder(r.Body).Decode(&condition); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}

	return UpdateConditionRequest{
		ID:   params.ByName("id"),
		Data: condition,
	}, nil
}

type UpdateConditionRequest struct {
	ID   string
	Data Condition
}

func makeInvoiceRoutes(prefix string, router *httprouter.Router, endpoints Set, options httptransport.ServerOption) {
	router.Handler(http.MethodPost, prefix+"/", httptransport.NewServer(
		endpoints.CreateInvoice,
		decodeCreateInvoiceRequest,
		httptransport.EncodeJSONResponse,
		options,
	))
}

func decodeCreateInvoiceRequest(ctx context.Context, r *http.Request) (any, error) {
	var invoice Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input and try again",
		)
	}
	return invoice, nil
}
