package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T)  {
	var h Server // NewServer 之类的方法
	//h := Server{}

	h.AddRoute(http.MethodGet, "/user", func(ctx Context) {
		fmt.Println("hello world!")
	})

	// 用法一：完全交给 http 包管理
	http.ListenAndServe(":8084", h)
	http.ListenAndServeTLS(":443", "", "" , h)

	// 用法二：自己手里管
	h.Start(":8085")
}