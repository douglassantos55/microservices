package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/inventory/pkg"
	"reconcip.com.br/microservices/inventory/proto"
)

func main() {
	supplierUrl := os.Getenv("SUPPLIER_SERVICE_URL")
	conn, err := grpc.Dial(supplierUrl+":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	validator := pkg.NewValidator([]pkg.ValidationRule{
		pkg.NewSupplierRule(conn),
	})

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
	endpoints = pkg.FetchSupplierEndpoints(endpoints, conn)

	var wg sync.WaitGroup
	wg.Add(3)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		user := os.Getenv("BROKER_USER")
		pass := os.Getenv("BROKER_PASSWORD")
		url := os.Getenv("BROKER_SERVICE_URL")

		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", user, pass, url))
		if err != nil {
			panic(err)
		}

		pkg.NewSubscriber(endpoints, conn)
	}(endpoints)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		authUrl := os.Getenv("AUTH_SERVICE_URL")
		conn, err := grpc.Dial(authUrl+":8080", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		authEndpoints := pkg.CreateAuthEndpoints(endpoints, conn)
		http.ListenAndServe(":80", pkg.NewHTTPHandler(authEndpoints))
	}(endpoints)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		defer grpcListener.Close()

		server := grpc.NewServer()
		grpcServer := pkg.NewGRPCServer(endpoints)
		proto.RegisterInventoryServer(server, grpcServer)

		if err := server.Serve(grpcListener); err != nil {
			panic(err)
		}
	}(endpoints)

	wg.Wait()
}
