package main

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/customer/pkg"
	"reconcip.com.br/microservices/customer/proto"
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

	endpoints := pkg.NewSet(svc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		server := grpc.NewServer()
		grpcServer := pkg.NewGRPCServer(endpoints)

		proto.RegisterCustomerServer(server, grpcServer)
		server.Serve(grpcListener)
	}(endpoints)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
		cc, err := grpc.Dial(authServiceUrl+":8080", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		defer cc.Close()

		verifyEndpoints := pkg.NewVerifySet(endpoints, cc)
		http.ListenAndServe(":80", pkg.NewHTTPHandler(verifyEndpoints))
	}(endpoints)

	wg.Wait()
}
