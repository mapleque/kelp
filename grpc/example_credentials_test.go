package grpc_test

import (
	"fmt"
	"golang.org/x/net/context"

	"github.com/mapleque/kelp/grpc"
	"github.com/mapleque/kelp/grpc/example"
)

// This example shows how to authorization with token
func Example_authorizeWithToken() {
	host := ":9999"
	token := "your_auth_token"
	go runServer(host, token)
	runClient(host, token)
	// Output:
	// message:"Hello kelp"
}

func runServer(host, token string) {
	gServer := grpc.New(
		grpc.TokenAuthorization(token),
	)
	server := &example.Server{}
	example.RegisterGreeterServer(gServer, server)
	grpc.Run(gServer, host)
}

func runClient(host, token string) {
	conn, err := grpc.DialWithToken(host, token)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	c := example.NewGreeterClient(conn)
	resp, err := c.SayHello(context.Background(), &example.HelloRequest{
		Name: "kelp",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
