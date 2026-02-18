package ex

/*
 * 这个文件是egine的内容，暂时先写这么多注释
 */
import "net/http"

type Engine struct {
	router      *Router
	middlewares []HandlerFunc
}

// 实例化引擎
func NewEngine() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

// 必须实现ServeHTTP方法
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Writer: w,
		Req:    r,
	}
	ctx.handlers = append(ctx.handlers, e.middlewares...)
	e.router.handle(ctx)
}

func (e *Engine) Use(middlewares ...HandlerFunc) {
	e.middlewares = append(e.middlewares, middlewares...)
}

// 启动一个http server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// 实现http GET请求
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute("GET", path, handler)
}
