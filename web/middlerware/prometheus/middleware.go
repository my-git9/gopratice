package prometheus

import (
    "github.com/prometheus/client_golang/prometheus"
    "gopratice/web"
    "strconv"
    "time"
)

type MiddlewareBuilder struct {
    Namespace string
    Subsystem string
    Name string
    Help string
}

func (m MiddlewareBuilder) Build() web.Middleware {
    vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
        Namespace: m.Namespace,
        Subsystem: m.Subsystem,
        Name:      m.Name,
        Help:      m.Help,
        Objectives: map[float64]float64{
            0.5:   0.05,
            0.75:  0.01,
            0.9:   0.01,
            0.99:  0.001,
        },
    }, []string{"pattern", "method", "status"})
    return func(next web.HandleFunc) web.HandleFunc {
        return func(ctx *web.Context) {
            startTime := time.Now()
            defer func() {
                duration := time.Now().Sub(startTime).Milliseconds()
                pattern := ctx.MatchedRoute
                if pattern == "" {
                    pattern = "unknown"
                }
                // WithLabelValues 是 Prometheus 提供的方法，用于为特定的标签组合（label values）选择一个指标实例
                vector.WithLabelValues(
                    pattern, // 请求路径模式（如 /api/v1/users）
                    ctx.Req.Method, // HTTP 请求方法（如 GET、POST）
                    strconv.Itoa(ctx.RespStatusCode), // HTTP 响应状态码（如 200、404）
                    ).Observe(float64(duration)) // 记录响应时间
            }()
            next(ctx)
        }
    }
}
