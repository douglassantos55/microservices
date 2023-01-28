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

	gateway := pkg.NewStripeGateway("sk_test_4eC39HqLyjWDarjtT1zdp7dc")

	svc := pkg.NewService(validator, repository, gateway)
	svc = pkg.NewLoggingService(svc, logger)

	endpoints := pkg.CreateEndpoints(svc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		authUrl := os.Getenv("AUTH_SERVICE_URL")
		ac, err := grpc.Dial(authUrl+":8080", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		customerUrl := os.Getenv("CUSTOMER_SERVICE_URL")
		cc, err := grpc.Dial(customerUrl+":8080", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		httpEndpoints := pkg.VerifyEndpoints(ac, endpoints)
		httpEndpoints = pkg.CustomerEndpoints(cc, endpoints)
		http.ListenAndServe(":80", pkg.NewHTTPHandler(httpEndpoints))
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
