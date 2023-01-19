package main

import (
	"net/http"
	"os"

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

	customerUrl := os.Getenv("CUSTOMER_SERVICE_URL")
	cc, err := grpc.Dial(customerUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	inventoryUrl := os.Getenv("INVENTORY_SERVICE_URL")
	ic, err := grpc.Dial(inventoryUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	validator := pkg.NewValidator([]pkg.ValidationRule{
		pkg.NewPaymentTypeRule(pc),
		pkg.NewPaymentMethodRule(pc),
		pkg.NewPaymentConditionRule(pc),
		pkg.NewCustomerRule(cc),
		pkg.NewEquipmentRule(ic),
	})

	deliveryUrl := os.Getenv("DELIVERY_SERVICE_URL")
	delivery, err := pkg.NewDeliveryService(deliveryUrl + ":8080")
	if err != nil {
		panic(err)
	}

	svc := pkg.NewService(validator, repository, delivery)

	endpoints := pkg.CreateEndpoints(svc)
	endpoints = pkg.WithPaymentMethodEndpoints(pc, endpoints)
	endpoints = pkg.WithPaymentTypeEndpoints(pc, endpoints)
	endpoints = pkg.WithPaymentConditionEndpoints(pc, endpoints)
	endpoints = pkg.WithCustomerEndpoints(cc, endpoints)

	http.ListenAndServe(":80", pkg.NewHTTPServer(endpoints))
}
