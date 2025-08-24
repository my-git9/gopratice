package web

// Middleware 函数式的责任链模式
// 函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

// AOP 方案在不同的框架，不同的语言中都有不同的叫法
// Middlewar, Hnadler, Chain, Filter, Filter-Chain, Interceptor, Wrapper
type MiddlewareV1 interface {
	Invoke(next HandleFunc) HandleFunc
}

type Interceptor interface {
	Before(ctx *Context)
	After(ctx *Context)
	Surround(ctx *Context)
}

type Chain []HandleFunc

type HandleFuncV1 func(ctx *Context) (next bool)

type ChainV1 struct {
	handlers []HandleFunc
}

func (c ChainV1) Run(ctx *Context) {
	for _, h := range c.handlers {
		next := h(ctx)
		if !next {
			return
		}
	}
}