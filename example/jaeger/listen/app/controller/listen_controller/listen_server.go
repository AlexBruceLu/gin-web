package listen_controller

import (
	"context"
	"fmt"
	"gin-web/example/jaeger/listen/app/proto/listen"
)

type ListenController struct{}

func (l *ListenController) ListenData(ctx context.Context, in *listen.Request) (*listen.Response, error) {
	return &listen.Response{Message: fmt.Sprintf("[%s]", in.Name)}, nil
}
