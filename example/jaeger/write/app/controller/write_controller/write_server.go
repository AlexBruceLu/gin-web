package write_controller

import (
	"fmt"
	"gin-web/example/jaeger/write/app/proto/listen"
	"gin-web/example/jaeger/write/app/proto/write"
	"gin-web/example/jaeger/write/app/util"
	"gin-web/example/jaeger/write/app/util/grpc_client"

	"golang.org/x/net/context"
)

type WriteController struct{}

func (s *WriteController) WriteData(ctx context.Context, in *write.Request) (*write.Response, error) {

	// 调用 gRPC 服务
	grpcListenClient := listen.NewListenClient(grpc_client.CreateServiceListenConn(ctx))
	resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})

	// 调用 HTTP 服务
	resHttpGet := ""
	_, err := util.HttpGet("http://localhost:9905/sing", ctx)
	if err == nil {
		resHttpGet = "[HttpGetOk]"
	}

	msg := "[" + fmt.Sprintf("%s", in.Name) + "-" +
		resListen.Message + "-" +
		resHttpGet +
		"]"
	return &write.Response{Message: msg}, nil
}
