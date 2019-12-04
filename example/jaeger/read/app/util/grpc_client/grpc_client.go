package grpc_client

import (
	"context"
	"gin-web/example/jaeger/read/app/util/jaeger_server"

	"github.com/smallnest/rpcx/log"
	"google.golang.org/grpc"
)

func CreateServerListenConn(ctx context.Context) *grpc.ClientConn {
	return createGRPCClient("127.0.0.1:9901", ctx)
}

func createGRPCClient(serviceAddress string, ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(jaeger_server.ClientInterceptor(jaeger_server.Tracer, ctx)))
	if err != nil {
		log.Error(serviceAddress, "grpc conn error", err)
	}
	return conn
}
