package main

import (
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"reconcip.com.br/microservices/delivery/pkg"
	"reconcip.com.br/microservices/delivery/proto"
)

func main() {
	svc := pkg.NewService([]pkg.Carrier{
		pkg.NewLocalCarrier(5, 7, pkg.NewMapeiaRouter(), pkg.NewMapeiaCoordinator()),
	})

	var wg sync.WaitGroup
	wg.Add(2)

	endpoints := pkg.CreateEndpoints(svc)

	go func(endpoints pkg.Set) {
		defer wg.Done()

		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		grpcServer := pkg.NewGRPCServer(endpoints)

		server := grpc.NewServer()
		proto.RegisterDeliveryServer(server, grpcServer)

		server.Serve(grpcListener)
	}(endpoints)

	go func(endpoints pkg.Set) {
		defer wg.Done()
		http.ListenAndServe(":80", pkg.NewHTTPServer(endpoints))
	}(endpoints)

	wg.Wait()
}
