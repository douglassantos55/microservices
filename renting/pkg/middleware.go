package pkg

import (
	"encoding/json"
	"time"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/go-kit/log"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(next Service, logger log.Logger) Service {
	return &loggingService{next, logger}
}

func (l *loggingService) CreateRent(data Rent) (rent *Rent, err error) {
	defer func() {
		l.logger.Log(
			"method", "CreateRent",
			"data", data,
			"rent", rent,
			"err", err,
		)
	}()
	return l.next.CreateRent(data)
}

func (l *loggingService) ListRents(page, perPage int64) (rents []*Rent, total int64, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListRents",
			"page", page,
			"perPage", perPage,
			"rents", rents,
			"total", total,
			"err", err,
		)
	}()
	return l.next.ListRents(page, perPage)
}

func (l *loggingService) UpdateRent(id string, data Rent) (rent *Rent, err error) {
	defer func() {
		l.logger.Log(
			"method", "UpdateRent",
			"id", id,
			"data", data,
			"rent", rent,
			"err", err,
		)
	}()
	return l.next.UpdateRent(id, data)
}

func (l *loggingService) DeleteRent(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "DeleteRent",
			"id", id,
			"err", err,
		)
	}()
	return l.next.DeleteRent(id)
}

func (l *loggingService) GetRent(id string) (rent *Rent, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetRent",
			"id", id,
			"rent", rent,
			"err", err,
		)
	}()
	return l.next.GetRent(id)
}

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
