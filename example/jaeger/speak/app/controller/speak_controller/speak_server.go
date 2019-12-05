package speak_controller

import (
	"context"
	"fmt"
	"gin-web/example/jaeger/speak/app/proto/listen"
	"gin-web/example/jaeger/speak/app/proto/speak"
	"gin-web/example/jaeger/speak/app/util"
	"gin-web/example/jaeger/speak/app/util/grpc_client"
)

type SpeakController struct {
}

func (s *SpeakController) SpeakData(ctx context.Context, in *speak.Request) (*speak.Response, error) {

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
	return &speak.Response{Message: msg}, nil
}
