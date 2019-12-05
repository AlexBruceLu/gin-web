package grpc_client

import (
	"context"
	"fmt"
	"gin-web/example/jaeger/speak/app/util/jaeger_server"

	"google.golang.org/grpc"
)

func CreateServiceListenConn(ctx context.Context) *grpc.ClientConn {
	return createGrpcClient("127.0.0.1:9901", ctx)
}

func createGrpcClient(serviceAddress string, ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithUnaryInterceptor(jaeger_server.ClientInterceptor(jaeger_server.Tracer, ctx)))
	if err != nil {
		fmt.Println(serviceAddress, "grpc conn err:", err)
	}
	return conn
}
