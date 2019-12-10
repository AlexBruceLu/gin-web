package grpc_log

import (
	"context"
	"fmt"
	"gin-web/app/config"
	utils "gin-web/example/jaeger/sing/app/util"
	"gin-web/example/jaeger/speak/app/util"

	"os"

	"github.com/smallnest/rpcx/log"
	"google.golang.org/grpc"
)

var grpcChannel = make(chan string, 100)

func handleGrpcChannel() {
	if f, err := os.OpenFile(config.AppAccessLogName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Warn(err)
	} else {
		for accessLog := range grpcChannel {
			f.WriteString(accessLog + "\n")
		}
	}
	return
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	go handleGrpcChannel()

	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := util.GetCurrentMilliUnix()

		err := invoker(ctx, method, req, reply, cc, opts...)

		endTime := util.GetCurrentMilliUnix()

		// 日志格式
		grpcLogMap := make(map[string]interface{})
		grpcLogMap["request_time"] = startTime
		grpcLogMap["request_data"] = req
		grpcLogMap["request_method"] = method
		grpcLogMap["response_data"] = reply
		grpcLogMap["response_error"] = err
		grpcLogMap["cost_time"] = fmt.Sprintf("%v ms", endTime-startTime)

		grpcLogJSON, _ := utils.JsonEncode(grpcLogMap)

		grpcChannel <- grpcLogJSON

		return err

	}
}
