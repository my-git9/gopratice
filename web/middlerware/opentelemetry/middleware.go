package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"gopratice/web"
)

const instrumentationName = "gitee.com/geektime-geekbang/geektime-go/web/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

/*
// NewMiddlewareBuilder 创建一个中间件构造器
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		Tracer: trace.
	}
}
*/

// Build 返回一个中间件
func (m MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		// 创建 tracer
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}

	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			// 尝试和客户端的 trace 结合
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			// 开始一个 span
			reqCtx, span := m.Tracer.Start(reqCtx, ctx.MatchedRoute)

			defer func() {
				// 这个只有执行完 next 才可能有值
				span.SetName(ctx.MatchedRoute)

				// 把响应码加上去
				span.SetAttributes(attribute.Int("http.status_code", ctx.RespStatusCode))
				span.End()
			}()

			// 设置 span 要记录的内容
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.path", ctx.Req.URL.Path))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))

			// 把 span 的上下文设置到 request 上
			ctx.Req = ctx.Req.WithContext(reqCtx)
			// 直接调用下一步
			next(ctx)

		}
	}
}
