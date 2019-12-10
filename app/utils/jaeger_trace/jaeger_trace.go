package jaeger_trace

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
	"google.golang.org/grpc/metadata"
)

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

func NewJaegerTracer(serviceName, jaegerHoastPort string) (opentracing.Tracer, io.Closer) {

	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHoastPort,
		},

		ServiceName: serviceName,
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot initialize jaeger", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func ClientInterceptor(tracer opentracing.Tracer, spanCtx opentracing.SpanContext) grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		span := opentracing.StartSpan(
			"call gRPC",
			opentracing.ChildOf(spanCtx),
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

		err := tracer.Inject(span.Context(), opentracing.TextMap, MDReaderWriter{md})
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
