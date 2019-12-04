package read_controller

import (
	"context"
	"fmt"
	"gin-web/example/jaeger/read/app/proto/listen"
	"gin-web/example/jaeger/read/app/proto/read"
	"gin-web/example/jaeger/read/app/util"
	"gin-web/example/jaeger/read/app/util/grpc_client"
)

type ReadController struct{}

func (i *ReadController) ReadData(ctx context.Context, in *read.Request) (*read.Response, error) {
	grpcListenCli := listen.NewListenClient(grpc_client.CreateServerListenConn(ctx))

	resListener, err := grpcListenCli.ListenData(context.Background(), &listen.Request{Name: "listen"})
	resHttp := ""

	_, err = util.HttpGet("http://localhost:9905/sing", ctx)
	if err == nil {
		resHttp = "[HttpGetOk]"
	}

	msg := "[" + fmt.Sprintf("%s", in.Name) + "-" + resListener.Message + "-" + resHttp + "]"

	return &read.Response{Message: msg}, nil

}
