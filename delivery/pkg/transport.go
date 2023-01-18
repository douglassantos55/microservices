package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
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

func (s *grpcServer) GetQuote(ctx context.Context, r *proto.GetQuoteRequest) (*proto.Quote, error) {
	_, reply, err := s.getQuote.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return reply.(*proto.Quote), nil
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

func NewHTTPServer(endpoints Set) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/", httptransport.NewServer(
		endpoints.GetQuotes,
		decodeGetQuotesRequest,
		httptransport.EncodeJSONResponse,
	))

	return router
}

func decodeGetQuotesRequest(ctx context.Context, r *http.Request) (any, error) {
	var req GetQuotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

type GetQuotesRequest struct {
	Origin string `json:"origin"`
	Dest   string `json:"dest"`
	Items  []Item `json:"items"`
}
