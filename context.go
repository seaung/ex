package ex

/*
 *  用于封装请求上下文
 */
import "net/http"

type HandlerFunc func(*Context)

// 请求上下文结构体
type Context struct {
	Writer   http.ResponseWriter
	Req      *http.Request
	handlers []HandlerFunc
	index    int
}

func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
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
