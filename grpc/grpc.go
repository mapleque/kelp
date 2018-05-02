package grpc

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"time"
)

// New build a grpc server with interceptors
func New(interceptors ...grpc.UnaryServerInterceptor) *grpc.Server {
	return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(interceptors...)))
}

// Run make the grpc server start serve on host
func Run(gServer *grpc.Server, host string) {
	reflection.Register(gServer)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("grpc service listen on", host)
	if err := gServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

// UnaryInterceptorChain wrap interceptors in one interceptor
func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		handlerChain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			handlerChain = buildHandler(interceptors[i], info, handlerChain)
		}
		return handlerChain(c, req)
	}
}

func buildHandler(interceptor grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, handlerChain grpc.UnaryHandler) grpc.UnaryHandler {
	return func(c context.Context, req interface{}) (interface{}, error) {
		return interceptor(c, req, info, handlerChain)
	}
}

// Recovery is an interceptor to recover when request deal panic
func Recovery(c context.Context, param interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = grpc.Errorf(codes.Internal, "panic error: %v", err)
			log.Error("[panic]", err)
			return
		}
	}()
	return handler(c, param)
}

// Logger is an interceptor to log request info
func Logger(c context.Context, param interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	resp, err = handler(c, param)

	end := time.Now()
	method := info.FullMethod

	latency := end.Sub(start)
	log.Info(
		"-", // remote ip
		end.Format("2006/01/02 15:04:05"),
		latency.Nanoseconds(),
		method,
		"-", // trace id
		"-", // uuid
		param,
		resp,
	)
	return
}
