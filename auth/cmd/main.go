package main

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/auth/pkg"
	"reconcip.com.br/microservices/auth/proto"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	svc := pkg.NewService(pkg.NewTokenGenerator())
	svc = pkg.NewLoggingService(svc, logger)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		defer func() {
			wg.Done()
			grpcListener.Close()
		}()

		server := grpc.NewServer()
		grpcServer := pkg.NewGRPCServer(svc)
		proto.RegisterAuthServer(server, grpcServer)

		if err := server.Serve(grpcListener); err != nil {
			logger.Log("failed", "serving grpc")
		}
	}()

	go func() {
		defer wg.Done()
		http.ListenAndServe(":80", pkg.NewHTTPHandler(svc))
	}()

	wg.Wait()
}
