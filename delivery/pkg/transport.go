package pkg

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	"reconcip.com.br/microservices/delivery/proto"
)

type grpcServer struct {
	proto.UnimplementedDeliveryServer
	getQuotes grpc.Handler
}

func NewGRPCServer(endpoints Set) proto.DeliveryServer {
	return &grpcServer{
		getQuotes: grpc.NewServer(
			endpoints.GetQuotes,
			decodeGetQuotesRequest,
			encodeGetQuotesResponse,
		),
	}
}

func decodeGetQuotesRequest(ctx context.Context, r any) (any, error) {
	req := r.(*proto.GetQuotesRequest)
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

	return GetQuotesRequest{
		Origin: req.GetOrigin(),
		Dest:   req.GetDestination(),
		Items:  items,
	}, nil
}

func encodeGetQuotesResponse(ctx context.Context, res any) (any, error) {
	reply := res.([]*Quote)
	quotes := make([]*proto.Quote, len(reply))

	for i, quote := range reply {
		quotes[i] = &proto.Quote{
			Carrier: quote.Carrier,
			Value:   quote.Value,
		}
	}

	return &proto.QuotesReply{Quotes: quotes}, nil
}

type GetQuotesRequest struct {
	Origin string
	Dest   string
	Items  []Item
}
