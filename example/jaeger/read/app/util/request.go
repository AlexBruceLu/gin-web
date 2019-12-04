package util

import (
	"context"
	"crypto/tls"
	"gin-web/example/jaeger/read/app/util/jaeger_server"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func HttpGet(url string, ctx context.Context) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "call http get",
		opentracing.Tag{Key: string(ext.Component), Value: "Http"},
		ext.SpanKindRPCClient,
	)
	span.Finish()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	injectErr := jaeger_server.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if injectErr != nil {
		log.Fatal("package read Tracer.Inject error", injectErr)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), nil

}
