//og:build e2e
package opentelemetry

import (
	"go.opentelemetry.io/otel"
    "testing"
	"gopratice/web"
	"gopratice/web/middlerware/accesslog"
	"time"

	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{
		Tracer: tracer,

	}
	server := web.NewHTTPServer(web.ServerWitMiddleware(builder.Build()))

	server.Get("/user", func(ctx *web.Context) {
		// 如果后面都是从这个 c 开始 start，就代表上包下，如果新启一个 ctx，则代表上下为并行 span
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second")
		time.Sleep(time.Second)
		c, third1 := tracer.Start(c, "third1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		c, third2 := tracer.Start(c, "third2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()

		c, first := tracer.Start(ctx.Req.Context(), "first_layer_1")
		defer first.End()
		ctx.Resp.Write([]byte("hello"))
	})

	server.Start(":8080")
}
