package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/supplier/pkg"
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

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	svc := pkg.NewService(validator, repository)
	svc = pkg.NewLoggingService(svc, logger)

	authServiceUrl := fmt.Sprintf("%s:8080", os.Getenv("AUTH_SERVICE_URL"))
	conn, err := grpc.Dial(authServiceUrl, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	endpoints := pkg.NewSet(svc)
	endpoints = pkg.NewVerifySet(endpoints, conn)

	http.ListenAndServe(":80", pkg.NewHTTPServer(endpoints))
}
