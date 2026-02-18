package ex

import "net/http"

// 封装请求上下文

type HandlerFunc func(*Context)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
}

// 响应http string
func (ctx *Context) String(code int, msg string) {
	ctx.Writer.WriteHeader(code)
	ctx.Writer.Write([]byte(msg))
}

// 获取URL参数
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}
