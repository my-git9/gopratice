//gobuild -race
package accesslog

import (
	"fmt"
	"gopratice/web"
	"testing"
)



func TestMiddlewareBuilderE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWitMiddleware(mdl))
	// 注册
	server.Get("/a/b/c", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello"))
	})
	server.Start(":8081")
}
