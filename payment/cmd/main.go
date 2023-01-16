package main

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/payment/pkg"
	"reconcip.com.br/microservices/payment/proto"
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

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	validator := pkg.NewValidator([]pkg.ValidationRule{
		pkg.NewPaymentTypeRule(repository),
	})
	svc := pkg.NewService(validator, repository)
	svc = pkg.NewLoggingService(svc, logger)

	endpoints := pkg.CreateEndpoints(svc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(endpoints pkg.Set) {
		defer wg.Done()
		http.ListenAndServe(":80", pkg.NewHTTPHandler(endpoints))
	}(endpoints)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		server := grpc.NewServer()
		grpcServer := pkg.NewGRPCServer(endpoints)

		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		proto.RegisterPaymentServer(server, grpcServer)
		server.Serve(grpcListener)
	}(endpoints)

	wg.Wait()
}
