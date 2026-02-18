package ex

/*
 * 这个文件是egine的内容，暂时先写这么多注释
 */
import "net/http"

type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

// 实例化引擎
func NewEngine() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{
		engine: e,
	}

	e.groups = []*RouterGroup{e.RouterGroup}
	return e
}

func DefaultEngine() *Engine {
	engine := NewEngine()
	engine.Use(Logger(), Recovery())
	return engine
}

// 必须实现ServeHTTP方法
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)
	e.router.handle(ctx)
}

func (e *Engine) addRoute(method, path string, handler HandlerFunc, middlewares []HandlerFunc) {
	handlers := append(middlewares, handler)
	e.router.addRoute(method, path, handlers)
}

func (e *Engine) Use(middlewares ...HandlerFunc) {
	e.middlewares = append(e.middlewares, middlewares...)
}

// 启动一个http server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// 实现http GET请求
func (e *Engine) GET(path string, handlers ...HandlerFunc) {
	e.RouterGroup.GET(path, handlers...)
}

// 实现http POST请求
func (e *Engine) POST(path string, handlers ...HandlerFunc) {
	e.RouterGroup.POST(path, handlers...)
}

// 实现http PUT请求
func (e *Engine) PUT(path string, handlers HandlerFunc) {
	e.RouterGroup.PUT(path, handlers)
}

// 实现http DELETE请求
func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.RouterGroup.DELETE(path, handler)
}
