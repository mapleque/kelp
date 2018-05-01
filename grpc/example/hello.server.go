package example

import (
	"golang.org/x/net/context"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct{}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello " + in.Name}, nil
}
