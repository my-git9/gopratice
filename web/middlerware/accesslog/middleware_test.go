package accesslog

import (
	"fmt"
	"gopratice/web"
	"net/http"
	"testing"
)



func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWitMiddleware(mdl))
	// 注册
	server.Post("/a/b/c", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello"))
	})
	request, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}

    server.ServeHTTP(nil, request)
}
