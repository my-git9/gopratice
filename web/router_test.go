package web

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T)  {
	// 1. 构造路由树
	// 2. 验证路由树
	testRoutes := []struct{
		method string
		path string
	}{
		{
			method: http.MethodGet,
			path: "/",
		},
		{
			method: http.MethodGet,
			path: "/user",
		},
		{
			method: http.MethodGet,
			path: "/user/home",
		},
		{
			method: http.MethodGet,
			path: "/order/detail",
		},
		{
			method: http.MethodPost,
			path: "/order/create",
		},
	}

	var mockHandler HandleFunc = func(ctx Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	// 断言路由树与预期的一致
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path: "user",
						handler: mockHandler,
						children:  map[string]*node{
							"home": &node{
								path: "home",
								handler: mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						children:  map[string]*node{
							"detail": &node{
								path: "detail",
								handler: mockHandler,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children:  map[string]*node{
							"create": &node{
								path: "create",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}
	// 不能通过 assert.Equal 比较，因为 HandleFunc 是不可比的
	// assert.Equal(t, wantRouter, r)
	msg, ok := wantRouter.equal(r)
		if !ok {
			fmt.Sprintf("error %s, %v", msg, ok)
		}

}

func (r router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("cat find correspond http method"), false
		}
		// v, dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}


func (n *node) equal(y *node) (string, bool)  {
	if n.path != y.path {
		return fmt.Sprintf("node path not match"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("subnode count not euqal"), false
	}

	// 比较 handler
	hHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if hHandler != yHandler {
		return fmt.Sprintf("handlers are not euqal"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("subnot %s not exist", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}