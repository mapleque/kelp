package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// TokenCredential implement grpc.credentials.PerRPCCredentials
type TokenCredential struct {
	token string
}

// DialWithToken used on dial server with token authorization
// which can be authorized by TokenAuthorization interceptor.
func DialWithToken(host, token string) (*grpc.ClientConn, error) {
	creds := &TokenCredential{token}
	return grpc.Dial(
		host,
		grpc.WithPerRPCCredentials(creds),
		grpc.WithInsecure(),
	)
}

func (this *TokenCredential) GetRequestMetadata(c context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": this.token,
	}, nil
}

func (this *TokenCredential) RequireTransportSecurity() bool {
	return false
}

// TokenAuthorization is an interceptor authorized request token from client.
func TokenAuthorization(token string) grpc.UnaryServerInterceptor {
	return func(c context.Context, param interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(c)
		if !ok {
			err = grpc.Errorf(codes.Internal, "authorization failed: token not found")
			log.Error("[authorization failed]", err)
			return
		}
		if len(md["authorization"]) < 1 || md["authorization"][0] != token {
			err = grpc.Errorf(codes.Internal, "authorization failed: invalid token %v", md)
			log.Error("[authorization failed]", err)
			return
		}
		return handler(c, param)
	}
}
