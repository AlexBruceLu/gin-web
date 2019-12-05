package jaeger_server

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

var Tracer opentracing.Tracer

func NewJaegerTracer(serviceName string, jaegerHostPort string) (opentracing.Tracer, io.Closer, error) {

	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: serviceName,
	}

	var closer io.Closer
	var err error

	Tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(Tracer)
	return Tracer, closer, err
}

func serverInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		var parentContext context.Context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		spanCtx, err := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Error("tracer extract from metadata error", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()
			parentContext = opentracing.ContextWithSpan(parentContext, span)
		}
		return handler(parentContext, req)

	}

}

func ServerOption(tracer opentracing.Tracer) grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor(tracer))
}

type MDReaderWriter struct {
	metadata.MD
}

// Set a key:value pair to the carrier. Multiple calls to Set() for the
// same key leads to undefined behavior.
//
// NOTE: The backing store for the TextMapWriter may contain data unrelated
// to SpanContext. As such, Inject() and Extract() implementations that
// call the TextMapWriter and TextMapReader interfaces must agree on a
// prefix or other convention to distinguish their own key:value pairs.
func (m MDReaderWriter) Set(key string, val string) {
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key], val)
}

// ForeachKey returns TextMap contents via repeated calls to the `handler`
// function. If any call to `handler` returns a non-nil error, ForeachKey
// terminates and returns that error.
//
// NOTE: The backing store for the TextMapReader may contain data unrelated
// to SpanContext. As such, Inject() and Extract() implementations that
// call the TextMapWriter and TextMapReader interfaces must agree on a
// prefix or other convention to distinguish their own key:value pairs.
//
// The "foreach" callback pattern reduces unnecessary copying in some cases
// and also allows implementations to hold locks while the map is read.
func (m MDReaderWriter) ForeachKey(handler func(key, val string) error) error {

	for k, vs := range m.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClientInterceptor grpc client
func ClientInterceptor(tracer opentracing.Tracer, ctx1 context.Context) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string,

		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		span, _ := opentracing.StartSpanFromContext(
			ctx1,
			"call gRPC",
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		mdWriter := MDReaderWriter{md}
		err := tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			span.LogFields(log.String("inject-error", err.Error()))
		}

		newCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.LogFields(log.String("call-error", err.Error()))
		}
		return err
	}
}
