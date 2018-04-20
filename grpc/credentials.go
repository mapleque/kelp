package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type TokenCredential struct {
	token string
}

func DialWithToken(host, token string) (*grpc.ClientConn, error) {
	//tokenCredential := &TokenCredential{token}
	return grpc.Dial(
		host,
		//grpc.WithPerRPCCredentials(tokenCredential),
		grpc.WithAuthority(token),
		grpc.WithInsecure(),
	)
}

func (this *TokenCredential) GetRequestMetadata(c context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authority": this.token,
	}, nil
}

func (this *TokenCredential) RequireTransportSecurity() bool {
	return true
}

func TokenAuthority(token string) grpc.UnaryServerInterceptor {
	return func(c context.Context, param interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(c)
		if !ok {
			err = grpc.Errorf(codes.Internal, "authority failed: token not found")
			log.Error("[authority failed]", err)
			return
		}
		if len(md[":authority"]) < 1 || md[":authority"][0] != token {
			err = grpc.Errorf(codes.Internal, "authority failed: invalid token %v", md)
			log.Error("[authority failed]", err)
			return
		}
		return handler(c, param)
	}
}
