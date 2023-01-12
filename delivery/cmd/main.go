package main

import (
	"net"

	"google.golang.org/grpc"
	"reconcip.com.br/microservices/delivery/pkg"
	"reconcip.com.br/microservices/delivery/proto"
)

func main() {
	svc := pkg.NewService([]pkg.Carrier{
		pkg.NewLocalCarrier(5, 7),
	})

	grpcListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	endpoints := pkg.CreateEndpoints(svc)
	grpcServer := pkg.NewGRPCServer(endpoints)

	server := grpc.NewServer()
	proto.RegisterDeliveryServer(server, grpcServer)

	server.Serve(grpcListener)
}
