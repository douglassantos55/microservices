package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/supplier/pkg"
	"reconcip.com.br/microservices/supplier/proto"
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

	endpoints := pkg.NewSet(svc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(endpoints pkg.Set) {
		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		defer grpcListener.Close()

		server := grpc.NewServer()
		grpcServer := pkg.NewGRPCServer(endpoints)
		proto.RegisterSupplierServiceServer(server, grpcServer)

		if err := server.Serve(grpcListener); err != nil {
			panic(err)
		}
	}(endpoints)

	go func(endpoints pkg.Set) {
		authServiceUrl := fmt.Sprintf("%s:8080", os.Getenv("AUTH_SERVICE_URL"))
		conn, err := grpc.Dial(authServiceUrl, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		httpEndpoints := pkg.NewVerifySet(endpoints, conn)
		http.ListenAndServe(":80", pkg.NewHTTPServer(httpEndpoints))
	}(endpoints)

	wg.Wait()
}
