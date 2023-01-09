package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/customer/pkg"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

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
	svc = pkg.NewLoggingService(svc, logger)

	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	cc, err := grpc.Dial(authServiceUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer cc.Close()

	endpoints := pkg.NewSet(svc)
	endpoints = pkg.NewVerifySet(endpoints, cc)
	http.ListenAndServe(":80", pkg.NewHTTPHandler(endpoints))
}
