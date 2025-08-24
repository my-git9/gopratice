package web

import (
	"log"
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 用来确保某个结构体一定实现了某个接口
// 确保一定实现了 Server 接口
var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(add string) error

	// 增加一些路由注册的功能
	// AddRoute 路由注册功能
	// method 是 HTTP 方法
	// path 是路由
	// handleFunc 是你的业务逻辑
	AddRoute(method string, path string, handleFunc HandleFunc)
	// 第二种方法：可以添加多个 handles 函数
	// AddRoute(method string, path string, handles... HandleFunc)
}

type HTTPServerOption func (s *HTTPServer)

type HTTPServer struct {
	// addr string 也可以创建的时候传递，而不是 Start 接收，这个都是可以的，具体看使用偏好
	*router

	mdls []Middleware

	tplEngine TemplateEngine
}

/*
func NewHTTPServer(mdls ...Middleware) *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
		mdls: mdls,
	}
}
 */

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithTemplateEngine(tplEngine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tplEngine
	}
}

func ServerWitMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	// 你的框架就在这里
	ctx := &Context{
		Req: request,
		Resp: writer,
		tplEngine: h.tplEngine,
	}

	// 最后一个是这个
	root := h.serve

	// 然后这里就是利用最后一个不断往前回溯组装链条
	// 从后往前
	// 
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}

	// 这里执行的时候就是从前往后了

	// 这里，最后一个步骤，就是把 RespData 和 RespStatusCode 刷新到响应里

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			// 这里就设置好了 RespData 和 RespStatusCode
			next(ctx)
			h.flashResp(ctx)
		}
	}

	root = m(root)

	root(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		log.Fatalln("write response error")
	}
}

func (h *HTTPServer) serve(ctx *Context)  {
	// 查找路由，并执行命中的业务逻辑
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.handler == nil {
		//ctx.Resp.WriteHeader(404)
		ctx.RespStatusCode = 404
		//ctx.Resp.Write([]byte("NOT FOUND"))
		ctx.RespData = []byte("NOT FOUND")
	}
	ctx.MatchedRoute = n.path
	n.handler(ctx)
}

/*
func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc)  {
	// TODO impement me
	panic("impement me")
}

 */

func (h *HTTPServer) Get(path string, handleFunc HandleFunc)  {
	h.AddRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc HandleFunc)  {
	h.AddRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServer) Start(addr string) error {
	// 也可以自己创建 Server
	// http.Server{}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 在这里，可以让用户注册所谓的 after start 回调
	// 比如往你的 admin 注册一些自己这个实例
	// 在这里执行一些业务所需的前置条件

	return http.Serve(l, h)
}