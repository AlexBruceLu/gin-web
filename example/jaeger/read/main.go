package main

import (
	"fmt"
	"gin-web/example/jaeger/read/app/controller/read_controller"
	"gin-web/example/jaeger/read/app/proto/read"
	"gin-web/example/jaeger/read/app/util/jaeger_server"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ServiceName     = "gRPC-Service-Read"
	ServiceHostPort = "0.0.0.0:9903"
	JaegerHostPort  = "127.0.0.1:6381"
)

func main() {
	var serviceOpts []grpc.ServerOption

	tracer, _, err := jaeger_server.NewJaegerTracer(ServiceName, JaegerHostPort)
	if err != nil {
		fmt.Printf("new tracer err: %+v\n", err)
		os.Exit(-1)
	}
	if tracer != nil {
		serviceOpts = append(serviceOpts, jaeger_server.ServerOption(tracer))
	}

	l, err := net.Listen("tcp", ServiceHostPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(serviceOpts...)

	// 服务注册
	read.RegisterReadServer(s, &read_controller.ReadController{})

	log.Println("Listen on " + ServiceHostPort)
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
