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
	next         Service
	reduceStock  endpoint.Endpoint
	processLater endpoint.Endpoint
}

func NewInventoryService(svc Service, reduceStock, processLater endpoint.Endpoint) Service {
	return &inventoryService{svc, reduceStock, processLater}
}

func (s *inventoryService) ListRents(page, perPage int64) ([]*Rent, int64, error) {
	return s.next.ListRents(page, perPage)
}

func (s *inventoryService) UpdateRent(id string, data Rent) (*Rent, error) {
	return s.next.UpdateRent(id, data)
}

func (s *inventoryService) CreateRent(data Rent) (*Rent, error) {
	rent, err := s.next.CreateRent(data)
	if err != nil {
		return nil, err
	}

	for _, item := range rent.Items {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if _, err := s.reduceStock(ctx, item); err != nil {
			s.processLater(ctx, item)
		}
	}

	return rent, nil
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
