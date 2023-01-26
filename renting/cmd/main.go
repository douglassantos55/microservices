package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/renting/pkg"
)

func main() {
	repository, err := pkg.NewMongoRepository(
		os.Getenv("MONGODB_URL"),
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_DATABASE"),
	)

	if err != nil {
		panic(err)
	}

	paymentUrl := os.Getenv("PAYMENT_SERVICE_URL")
	pc, err := grpc.Dial(paymentUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	customerUrl := os.Getenv("CUSTOMER_SERVICE_URL")
	cc, err := grpc.Dial(customerUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	inventoryUrl := os.Getenv("INVENTORY_SERVICE_URL")
	ic, err := grpc.Dial(inventoryUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer ic.Close()

	validator := pkg.NewValidator([]pkg.ValidationRule{
		pkg.NewPaymentTypeRule(pc),
		pkg.NewPaymentMethodRule(pc),
		pkg.NewPaymentConditionRule(pc),
		pkg.NewCustomerRule(cc),
	})

	deliveryUrl := os.Getenv("DELIVERY_SERVICE_URL")
	dc, err := grpc.Dial(deliveryUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer dc.Close()

	delivery := pkg.NewGRPCDeliveryService(dc)

	brokerUrl := os.Getenv("BROKER_SERVICE_URL")
	brokerUser := os.Getenv("BROKER_USER")
	brokerPass := os.Getenv("BROKER_PASSWORD")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", brokerUser, brokerPass, brokerUrl))
	if err != nil {
		panic(err)
	}

	inventory := pkg.NewInventoryService(
		pkg.ReduceStockEndpoint(ic),
		pkg.RestoreStockEndpoint(ic),
		pkg.ProcessLaterEndpoint(conn),
	)

	svc := pkg.NewService(validator, repository, delivery, inventory)

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	svc = pkg.NewLoggingService(svc, logger)

	endpoints := pkg.CreateEndpoints(svc)
	endpoints = pkg.WithEquipmentEndpoints(ic, endpoints)
	endpoints = pkg.WithPaymentMethodEndpoints(pc, endpoints)
	endpoints = pkg.WithPaymentTypeEndpoints(pc, endpoints)
	endpoints = pkg.WithPaymentConditionEndpoints(pc, endpoints)
	endpoints = pkg.WithCustomerEndpoints(cc, endpoints)

	http.ListenAndServe(":80", pkg.NewHTTPServer(endpoints))
}
