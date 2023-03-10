package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/auth/jwt"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
	"reconcip.com.br/microservices/inventory/proto"
)

func NewHTTPHandler(endpoints Set) http.Handler {
	router := httprouter.New()

	options := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		endpoints.Create,
		decodeCreateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodGet, "/", httptransport.NewServer(
		endpoints.List,
		decodeListRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodPut, "/:id", httptransport.NewServer(
		endpoints.Update,
		decodeUpdateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Handler(http.MethodDelete, "/:id", httptransport.NewServer(
		endpoints.Delete,
		URLParamDecoder("id"),
		encodeDeleteResponse,
		options...,
	))

	router.Handler(http.MethodGet, "/:id", httptransport.NewServer(
		endpoints.Get,
		URLParamDecoder("id"),
		httptransport.EncodeJSONResponse,
		options...,
	))

	return router
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	var equipment Equipment
	if err := json.NewDecoder(r.Body).Decode(&equipment); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data and try again",
		)
	}
	return equipment, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	params := r.URL.Query()
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(params.Get("per_page"))
	if err != nil {
		perPage = 50
	}

	return Pagination{page - 1, perPage}, nil
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	var equipment Equipment
	if err := json.NewDecoder(r.Body).Decode(&equipment); err != nil {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid input data",
			"verify your input data and try again",
		)
	}

	return UpdateRequest{
		ID:   params.ByName("id"),
		Data: equipment,
	}, nil
}

type UpdateRequest struct {
	ID   string    `json:"id"`
	Data Equipment `json:"data"`
}

func URLParamDecoder(param string) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (any, error) {
		params := httprouter.ParamsFromContext(r.Context())
		return params.ByName(param), nil
	}
}

func encodeDeleteResponse(ctx context.Context, w http.ResponseWriter, r any) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

type grpcServer struct {
	proto.UnimplementedInventoryServer
	reduceStock  grpc.Handler
	getEquipment grpc.Handler
	restoreStock grpc.Handler
}

func NewGRPCServer(endpoints Set) proto.InventoryServer {
	return &grpcServer{
		reduceStock: grpc.NewServer(
			endpoints.ReduceStock,
			decodeReduceStockRequest,
			NopGRPCEncoder,
		),
		getEquipment: grpc.NewServer(
			endpoints.Get,
			decodeGetRequest,
			encodeEquipmentResponse,
		),
		restoreStock: grpc.NewServer(
			endpoints.RestoreStock,
			decodeRestoreStockRequest,
			NopGRPCEncoder,
		),
	}
}

func (s *grpcServer) GetEquipment(ctx context.Context, r *proto.GetRequest) (*proto.Equipment, error) {
	_, reply, err := s.getEquipment.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return reply.(*proto.Equipment), nil
}

func (s *grpcServer) ReduceStock(ctx context.Context, req *proto.ReduceStockRequest) (*proto.ReduceStockReply, error) {
	_, _, err := s.reduceStock.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.ReduceStockReply{Err: err.Error()}, nil
	}
	return &proto.ReduceStockReply{}, nil
}

func (s *grpcServer) RestoreStock(ctx context.Context, req *proto.RestoreStockRequest) (*proto.RestoreStockReply, error) {
	_, _, err := s.restoreStock.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return &proto.RestoreStockReply{}, nil
}

func decodeReduceStockRequest(ctx context.Context, req any) (any, error) {
	item := req.(*proto.ReduceStockRequest)

	return ReduceStockRequest{
		Equip: item.GetId(),
		Qty:   item.GetQty(),
	}, nil
}

func decodeRestoreStockRequest(ctx context.Context, r any) (any, error) {
	item := r.(*proto.RestoreStockRequest)

	return RestoreStockRequest{
		Equip: item.GetId(),
		Qty:   item.GetQty(),
	}, nil
}

func NopGRPCEncoder(ctx context.Context, res any) (any, error) {
	return nil, nil
}

type ReduceStockRequest struct {
	Equip string `json:"equip_id"`
	Qty   int64  `json:"qty"`
}

type RestoreStockRequest struct {
	Equip string `json:"equip_id"`
	Qty   int64  `json:"qty"`
}

func decodeGetRequest(ctx context.Context, r any) (any, error) {
	req := r.(*proto.GetRequest)
	return req.GetId(), nil
}

func encodeEquipmentResponse(ctx context.Context, r any) (any, error) {
	equipment := r.(*Equipment)
	supplier := &proto.Supplier{}

	if equipment.Supplier != nil {
		supplier.Id = equipment.Supplier.ID
		supplier.SocialName = equipment.Supplier.SocialName
		supplier.LegalName = equipment.Supplier.LegalName
		supplier.Email = equipment.Supplier.Email
		supplier.Website = equipment.Supplier.Website
		supplier.Cnpj = equipment.Supplier.Cnpj
		supplier.InscEst = equipment.Supplier.InscEst
		supplier.Phone = equipment.Supplier.Phone
	}

	rentingValues := make([]*proto.RentingValue, len(equipment.RentingValues))
	for i, value := range equipment.RentingValues {
		rentingValues[i] = &proto.RentingValue{
			Value: value.Value,
			Period: &proto.Period{
				Id:      value.PeriodID,
				Name:    value.PeriodID,
				QtyDays: 55,
			},
		}
	}

	return &proto.Equipment{
		Id:             equipment.ID,
		Description:    equipment.Description,
		Stock:          int64(equipment.Stock),
		EffectiveStock: int64(equipment.EffectiveStock),
		Weight:         equipment.Weight,
		UnitValue:      equipment.UnitValue,
		PurchaseValue:  equipment.PurchaseValue,
		ReplaceValue:   equipment.ReplaceValue,
		MinQty:         int64(equipment.MinQty),
		Supplier:       supplier,
		RentingValues:  rentingValues,
	}, nil
}

func NewSubscriber(endpoints Set, conn *amqp.Connection) {
	subscriber := amqptransport.NewSubscriber(
		endpoints.ReduceStock,
		decodeReduceStockAMQPRequest,
		amqptransport.EncodeNopResponse,
	)

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()

	if err := channel.ExchangeDeclare("inventory", "direct", true, false, false, false, nil); err != nil {
		panic(err)
	}

	queue, err := channel.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	if err := channel.QueueBind(queue.Name, "stock.reduce", "inventory", false, nil); err != nil {
		panic(err)
	}

	handler := subscriber.ServeDelivery(channel)
	messages, err := channel.Consume(queue.Name, "", false, true, false, false, nil)

	go func(<-chan amqp.Delivery) {
		for message := range messages {
			handler(&message)
		}
	}(messages)

	var forever chan any
	<-forever
}

func decodeReduceStockAMQPRequest(ctx context.Context, d *amqp.Delivery) (any, error) {
	var item struct {
		EquipmentID string `json:"equipment_id"`
		Qty         int    `json:"qty"`
	}

	if err := json.Unmarshal(d.Body, &item); err != nil {
		return nil, err
	}

	return ReduceStockRequest{
		Equip: item.EquipmentID,
		Qty:   int64(item.Qty),
	}, nil
}
