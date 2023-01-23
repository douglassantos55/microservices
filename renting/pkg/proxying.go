package pkg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/renting/proto"
)

func WithPaymentTypeEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withPaymentType := withPaymentTypeMiddleware(cc)
	return Set{
		Create: withPaymentType(endpoints.Create),
		List:   withPaymentType(endpoints.List),
		Update: withPaymentType(endpoints.Update),
		Delete: endpoints.Delete,
	}
}

func withPaymentTypeMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getPaymentType := getPaymentTypeEndpoint(cc)

	appendType := func(ctx context.Context, rent *Rent) {
		paymentType, err := getPaymentType(ctx, rent.PaymentTypeID)
		if err == nil {
			rent.PaymentType = paymentType.(*PaymentType)
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				appendType(ctx, rent)
			}

			if result, ok := res.(ListResult); ok {
				for _, rent := range result.Items {
					appendType(ctx, rent.(*Rent))
				}
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
		Create: withPaymentMethod(endpoints.Create),
		List:   withPaymentMethod(endpoints.List),
		Update: withPaymentMethod(endpoints.Update),
		Delete: endpoints.Delete,
	}
}

func withPaymentMethodMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getPaymentMethod := getPaymentMethodEndpoint(cc)

	appendMethod := func(ctx context.Context, rent *Rent) {
		method, err := getPaymentMethod(ctx, rent.PaymentMethodID)
		if err == nil {
			rent.PaymentMethod = method.(*PaymentMethod)
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				appendMethod(ctx, rent)
			}

			if result, ok := res.(ListResult); ok {
				for _, rent := range result.Items {
					appendMethod(ctx, rent.(*Rent))
				}
			}

			return res, nil
		}
	}
}

func getPaymentMethodEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
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

func WithPaymentConditionEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withPaymentCondition := withPaymentConditionMiddleware(cc)

	return Set{
		Create: withPaymentCondition(endpoints.Create),
		List:   withPaymentCondition(endpoints.List),
		Update: withPaymentCondition(endpoints.Update),
		Delete: endpoints.Delete,
	}
}

func withPaymentConditionMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getPaymentCondition := getPaymentConditionEndpoint(cc)

	appendCondition := func(ctx context.Context, rent *Rent) {
		condition, err := getPaymentCondition(ctx, rent.PaymentConditionID)
		if err == nil {
			rent.PaymentCondition = condition.(*PaymentCondition)
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				appendCondition(ctx, rent)
			}

			if result, ok := res.(ListResult); ok {
				for _, rent := range result.Items {
					appendCondition(ctx, rent.(*Rent))
				}
			}

			return res, nil
		}
	}
}

func getPaymentConditionEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Payment",
		"GetCondition",
		encodeRequest,
		decodePaymentCondition,
		&proto.ConditionReply{},
	).Endpoint()
}

func decodePaymentCondition(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.ConditionReply)

	condition := reply.GetCondition()
	paymentType := condition.GetPaymentType()

	return &PaymentCondition{
		ID:        condition.GetId(),
		Name:      condition.GetName(),
		Increment: condition.GetIncrement(),
		PaymentType: &PaymentType{
			ID:   paymentType.GetId(),
			Name: paymentType.GetName(),
		},
		Installments: condition.GetInstallments(),
	}, nil
}

func WithCustomerEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withCustomer := withCustomerMiddleware(cc)

	return Set{
		Create: withCustomer(endpoints.Create),
		List:   withCustomer(endpoints.List),
		Update: withCustomer(endpoints.Update),
		Delete: endpoints.Delete,
	}
}

func withCustomerMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getCustomer := getCustomerEndpoint(cc)

	appendCustomer := func(ctx context.Context, rent *Rent) {
		customer, err := getCustomer(ctx, rent.CustomerID)
		if err == nil {
			rent.Customer = customer.(*Customer)
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if rent, ok := res.(*Rent); ok {
				appendCustomer(ctx, rent)
			}

			if result, ok := res.(ListResult); ok {
				for _, rent := range result.Items {
					appendCustomer(ctx, rent.(*Rent))
				}
			}

			return res, nil
		}
	}
}

func getCustomerEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Customer",
		"Get",
		encodeRequest,
		decodeCustomer,
		&proto.Customer{},
	).Endpoint()
}

func decodeCustomer(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.Customer)

	return &Customer{
		ID:        reply.GetId(),
		Name:      reply.GetName(),
		Email:     reply.GetEmail(),
		CpfCnpj:   reply.GetCpfCnpj(),
		RgInscEst: reply.GetRgInscEst(),
		Phone:     reply.GetPhone(),
		Cellphone: reply.GetCellphone(),
	}, nil
}

type grpcDeliveryService struct {
	conn *grpc.ClientConn
}

func NewGRPCDeliveryService(cc *grpc.ClientConn) DeliveryService {
	return &grpcDeliveryService{cc}
}

func (s *grpcDeliveryService) GetQuote(origin, dest, carrier string, items []*Item) (*Quote, error) {
	quoteItems := make([]QuoteItem, len(items))
	for i, item := range items {
		quoteItems[i] = QuoteItem{Qty: int64(item.Qty)}
	}

	getQuote := getQuoteEndpoint(s.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	quote, err := getQuote(ctx, QuoteRequest{
		Origin:  origin,
		Dest:    dest,
		Carrier: carrier,
		Items:   quoteItems,
	})

	if err != nil {
		return nil, err
	}

	return quote.(*Quote), nil
}

func getQuoteEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Delivery",
		"GetQuote",
		encodeQuoteRequest,
		decodeQuoteResponse,
		&proto.Quote{},
	).Endpoint()
}

func encodeQuoteRequest(ctx context.Context, r any) (any, error) {
	req := r.(QuoteRequest)
	items := make([]*proto.Item, len(req.Items))

	for i, item := range req.Items {
		items[i] = &proto.Item{
			Qty:    item.Qty,
			Weight: item.Weight,
			Width:  item.Width,
			Height: item.Height,
			Depth:  item.Depth,
		}
	}

	return &proto.GetQuoteRequest{
		Origin:      req.Origin,
		Destination: req.Dest,
		Carrier:     req.Carrier,
		Items:       items,
	}, nil
}

func decodeQuoteResponse(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.Quote)

	return &Quote{
		Carrier: reply.GetCarrier(),
		Value:   reply.GetValue(),
	}, nil
}

type QuoteRequest struct {
	Origin  string
	Dest    string
	Carrier string
	Items   []QuoteItem
}

type QuoteItem struct {
	Qty    int64
	Weight float64
	Width  float64
	Height float64
	Depth  float64
}

func WithEquipmentEndpoints(cc *grpc.ClientConn, endpoints Set) Set {
	withEquipment := withEquipmentMiddleware(cc)

	return Set{
		Create: withEquipment(endpoints.Create),
		List:   withEquipment(endpoints.List),
		Update: withEquipment(endpoints.Update),
		Delete: endpoints.Delete,
	}
}

func withEquipmentMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	getEquipment := getEquipmentEndpoint(cc)

	appendEquipment := func(ctx context.Context, index int, item *Item) error {
		equipment, err := getEquipment(ctx, item.EquipmentID)
		if err != nil {
			return NewError(
				http.StatusBadRequest,
				"equipment not found",
				fmt.Sprintf("Items[%d] equipment not found", index),
			)
		}
		item.Equipment = equipment.(*Equipment)
		return nil
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			if rent, ok := r.(Rent); ok {
				for i, item := range rent.Items {
					if err := appendEquipment(ctx, i, item); err != nil {
						return nil, err
					}
				}
				return next(ctx, rent)
			}

			if result, ok := r.(ListResult); ok {
				for _, item := range result.Items {
					for i, item := range item.(*Rent).Items {
						if err := appendEquipment(ctx, i, item); err != nil {
							return nil, err
						}
					}
				}
				return next(ctx, result)
			}

			if req, ok := r.(UpdateRequest); ok {
				for i, item := range req.Data.Items {
					if err := appendEquipment(ctx, i, item); err != nil {
						return nil, err
					}
				}
				return next(ctx, req)
			}

			return next(ctx, r)
		}
	}
}

func getEquipmentEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Inventory",
		"GetEquipment",
		encodeRequest,
		decodeEquipment,
		&proto.Equipment{},
	).Endpoint()
}

func decodeEquipment(ctx context.Context, r any) (any, error) {
	equipment := r.(*proto.Equipment)
	rentingValues := make([]*RentingValue, len(equipment.GetRentingValues()))

	for i, value := range equipment.GetRentingValues() {
		rentingValues[i] = &RentingValue{
			Value:    value.GetValue(),
			PeriodID: value.GetPeriod().GetId(),
			Period: &Period{
				ID:      value.GetPeriod().GetId(),
				Name:    value.GetPeriod().GetName(),
				QtyDays: value.GetPeriod().GetQtyDays(),
			},
		}
	}

	return &Equipment{
		ID:             equipment.GetId(),
		Description:    equipment.GetDescription(),
		Weight:         equipment.GetWeight(),
		UnitValue:      equipment.GetUnitValue(),
		EffectiveStock: int(equipment.GetEffectiveStock()),
		RentingValues:  rentingValues,
	}, nil
}

func RestoreStockEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Inventory",
		"RestoreStock",
		encodeRestoreStockRequest,
		NopGRPCDecoder,
		&proto.RestoreStockReply{},
	).Endpoint()
}

func ReduceStockEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Inventory",
		"ReduceStock",
		encodeReduceStockRequest,
		NopGRPCDecoder,
		&proto.ReduceStockReply{},
	).Endpoint()
}

func encodeReduceStockRequest(ctx context.Context, r any) (any, error) {
	item := r.(*Item)

	return &proto.ReduceStockRequest{
		Id:  item.EquipmentID,
		Qty: int64(item.Qty),
	}, nil
}

func encodeRestoreStockRequest(ctx context.Context, r any) (any, error) {
	item := r.(*Item)

	return &proto.RestoreStockRequest{
		Id:  item.EquipmentID,
		Qty: int64(item.Qty),
	}, nil
}

func NopGRPCDecoder(ctx context.Context, r any) (any, error) {
	return nil, nil
}
