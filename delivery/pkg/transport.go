package pkg

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"reconcip.com.br/microservices/delivery/proto"
)

type grpcServer struct {
	proto.UnimplementedDeliveryServer
	getQuote grpctransport.Handler
}

func NewGRPCServer(endpoints Set) proto.DeliveryServer {
	return &grpcServer{
		getQuote: grpctransport.NewServer(
			endpoints.GetQuote,
			decodeGetQuoteRequest,
			encodeGetQuoteResponse,
		),
	}
}

func decodeGetQuoteRequest(ctx context.Context, r any) (any, error) {
	req := r.(*proto.GetQuoteRequest)
	items := make([]Item, len(req.Items))

	for i, item := range req.GetItems() {
		items[i] = Item{
			Qty:    int(item.GetQty()),
			Weight: item.GetWeight(),
			Width:  item.GetWidth(),
			Height: item.GetHeight(),
			Depth:  item.GetDepth(),
		}
	}

	return GetQuoteRequest{
		Origin:  req.GetOrigin(),
		Dest:    req.GetDestination(),
		Carrier: req.GetCarrier(),
		Items:   items,
	}, nil
}

func encodeGetQuoteResponse(ctx context.Context, res any) (any, error) {
	reply := res.(*Quote)

	return &proto.Quote{
		Carrier: reply.Carrier,
		Value:   reply.Value,
	}, nil
}

type GetQuoteRequest struct {
	Origin  string
	Dest    string
	Carrier string
	Items   []Item
}

}

type GetQuotesRequest struct {
	Origin string
	Dest   string
	Items  []Item
}
