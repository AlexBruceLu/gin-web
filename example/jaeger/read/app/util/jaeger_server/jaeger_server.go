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

func ServerOption(tracer opentracing.Tracer) grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor(tracer))
}

var Tracer opentracing.Tracer

func NewJaegerTracer(serviceName, jaegerHostPort string) (opentracing.Tracer, io.Closer, error) {
	var closer io.Closer
	var err error

	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: serviceName,
	}
	Tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot initialize jaeger:err ", err))
	}

	opentracing.SetGlobalTracer(Tracer)
	return Tracer, closer, err

}

func ClientInterceptor(tracer opentracing.Tracer, ctx context.Context) grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoke grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		span, _ := opentracing.StartSpanFromContext(
			ctx,
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
		var err error
		if err = tracer.Inject(span.Context(), opentracing.TextMap, mdWriter); err != nil {
			span.LogFields(log.String("inject error", err.Error()))
		}

		newCtx := metadata.NewOutgoingContext(ctx, md)
		if err = invoke(newCtx, method, req, reply, cc, opts...); err != nil {
			span.LogFields(log.String("invoke error", err.Error()))
		}

		return err
	}

}

func serverInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (reps interface{}, err error) {

		var parentContext context.Context

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		spanContext, err := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Error("extract from metadata error: ", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()

			parentContext = opentracing.ContextWithSpan(ctx, span)
		}

		return handler(parentContext, req)
	}
}

type MDReaderWriter struct {
	metadata.MD
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
