package prometheus

import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "gopratice/web"
    "math/rand"
    "net/http"
    "testing"
    "time"
)

func TestMiddlewareBuilder_Build(t *testing.T)  {
    builder := MiddlewareBuilder{
        Namespace: "geekbang",
        Subsystem: "web",
        Name: "http_response",
        Help: "prometheus middleware",
    }
    server := web.NewHTTPServer(web.ServerWitMiddleware(builder.Build()))

    server.Get("/user", func(ctx *web.Context) {
        ctx.RespJSON(200, map[string]string{
            val := rand.Intn(1000)
            time.Sleep(time.Duration(val) * time.Millisecond)
            "name": "geekbang",
        })
    })

    go func() {
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":8082", nil)
    }()

    server.Start(":8080")
}

type User struct {
    Name string `json:"name"`
}