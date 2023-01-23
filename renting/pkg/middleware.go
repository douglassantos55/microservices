package pkg

import (
	"encoding/json"
	"time"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type inventoryService struct {
	reduceStock  endpoint.Endpoint
	restoreStock endpoint.Endpoint
	processLater endpoint.Endpoint
}

func (s *inventoryService) ReduceStock(items []*Item) {
	for _, item := range items {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if _, err := s.reduceStock(ctx, item); err != nil {
			s.processLater(ctx, item)
		}
	}
}

func (s *inventoryService) RestoreStock(items []*Item) {
	for _, item := range items {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		s.restoreStock(ctx, item)
	}
}

func NewInventoryService(reduceStock, restoreStock, processLater endpoint.Endpoint) InventoryService {
	return &inventoryService{reduceStock, restoreStock, processLater}
}

func ProcessLaterEndpoint(conn *amqp.Connection) endpoint.Endpoint {
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	replyQueue, err := channel.QueueDeclare("", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	return amqptransport.NewPublisher(
		channel,
		&replyQueue,
		encodeAMQPRequest,
		decodeAMQPResponse,
		amqptransport.PublisherBefore(
			amqptransport.SetPublishKey("stock.reduce"),
			amqptransport.SetPublishExchange("inventory"),
			amqptransport.SetContentType("application/json"),
		),
		amqptransport.PublisherDeliverer(amqptransport.SendAndForgetDeliverer),
	).Endpoint()
}

func encodeAMQPRequest(ctx context.Context, p *amqp.Publishing, r any) error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	p.Body = body
	return nil
}

func decodeAMQPResponse(ctx context.Context, d *amqp.Delivery) (any, error) {
	return nil, nil
}
