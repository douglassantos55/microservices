package main

import (
	"net/http"
	"os"

	"reconcip.com.br/microservices/payment/pkg"
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

	svc := pkg.NewService(pkg.NewValidator(), repository)
	endpoints := pkg.CreateEndpoints(svc)

	http.ListenAndServe(":80", pkg.NewHTTPHandler(endpoints))
}
