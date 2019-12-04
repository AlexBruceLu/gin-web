package main

import (
	"fmt"
	"gin-web/example/jaeger/listen/app/controller/listen_controller"
	"gin-web/example/jaeger/listen/app/proto/listen"
	"gin-web/example/jaeger/listen/app/util/jaeger_server"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ServiceName     = "gRPC-Service_Listen"
	ServiceHostPort = "0.0.0.0:9901"
	JaegerHostPort  = "127.0.0.1:6831"
)

func main() {
	var serviceOpts []grpc.ServerOption

	tracer, _, err := jaeger_server.NewJaegerTracer(ServiceName, JaegerHostPort)
	if err != nil {
		fmt.Println("Create Jaeger Tracer error: ", err)
		os.Exit(1)
	}
	if tracer != nil {
		serviceOpts = append(serviceOpts, jaeger_server.ServerOption(tracer))
	}

	l, err := net.Listen("tcp", ServiceHostPort)
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s := grpc.NewServer(serviceOpts...)

	listen.RegisterListenServer(s, &listen_controller.ListenController{})
	log.Println("Listen on", ServiceHostPort)
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatal("Failure to serve", err)

	}
}
