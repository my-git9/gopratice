package errhdl

import "gopratice/web"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func NewMiddlewareBuilder() MiddlewareBuilder {
	return MiddlewareBuilder{
		resp: make(map[int][]byte),
	}
}

func (m MiddlewareBuilder) Build() web.Middleware{
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				// 篡改结果
				ctx.RespData = resp
			}
		}
	}
}
