package request

import (
	"crypto/tls"
	"gin-web/app/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func HttpGet(url string, c *gin.Context) (string, error) {

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

	if config.JaegerOpen == 1 {

		tracer, _ := c.Get("Tracer")
		parentSpanContext, _ := c.Get("ParentSpanContext")

		span := opentracing.StartSpan(
			"call Http Get",
			opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
			opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			ext.SpanKindRPCClient,
		)

		span.Finish()

		injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if injectErr != nil {
			log.Fatal(err, "couldn't inject headers")
		}

	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	contect, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(contect), nil
}
