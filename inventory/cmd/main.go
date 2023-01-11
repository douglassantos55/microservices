package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/inventory/pkg"
)

func main() {
	validator := pkg.NewValidator()

	repository, err := pkg.NewMongoRepository(
		os.Getenv("MONGODB_URL"),
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_DATABASE"),
	)
	if err != nil {
		panic(err)
	}

	supplierUrl := os.Getenv("SUPPLIER_URL")
	conn, err := grpc.Dial(supplierUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	svc := pkg.NewService(validator, repository)
	svc = pkg.NewLoggingService(svc, logger)

	endpoints := pkg.NewSet(svc)
	endpoints = pkg.FetchSupplierEndpoints(endpoints, conn)

	http.ListenAndServe(":80", pkg.NewHTTPHandler(endpoints))
}
