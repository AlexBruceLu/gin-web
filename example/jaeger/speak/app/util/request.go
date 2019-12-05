package util

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func HttpGet(url string, ctx context.Context) (string, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		"call http get",
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
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
}
