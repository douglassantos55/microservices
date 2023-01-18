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

	validator := pkg.NewValidator([]pkg.ValidationRule{
		pkg.NewPaymentTypeRule(pc),
		pkg.NewPaymentMethodRule(pc),
		pkg.NewPaymentConditionRule(pc),
		pkg.NewCustomerRule(cc),
	})

	svc := pkg.NewService(validator, repository)

	endpoints := pkg.CreateEndpoints(svc)
	endpoints = pkg.WithPaymentMethodEndpoints(cc, endpoints)
	endpoints = pkg.WithPaymentTypeEndpoints(cc, endpoints)
	endpoints = pkg.WithPaymentConditionEndpoints(cc, endpoints)

	http.ListenAndServe(":80", pkg.NewHTTPServer(endpoints))
}
