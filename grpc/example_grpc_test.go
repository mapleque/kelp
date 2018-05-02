package grpc_test

import (
	"github.com/mapleque/kelp/grpc"
	"github.com/mapleque/kelp/grpc/example"
)

// This example shows how to run a grpc server
func Example_grpc() {
	host := ":9999"
	gServer := grpc.New(
		grpc.Recovery,
		grpc.Logger,
	)

	server := &example.Server{}
	example.RegisterGreeterServer(gServer, server)
	grpc.Run(gServer, host)
}
