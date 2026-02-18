package ex

/*
 * 这个文件是egine的内容，暂时先写这么多注释
 */
import "net/http"

type Engine struct {
	router *Router
}

// 实例化引擎
func NewEngine() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

// 必须实现ServeHTTP方法
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.handle(w, r)
}

// 启动一个http server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// 实现http GET请求
func (e *Engine) GET(path string, handler http.HandlerFunc) {
	e.router.addRoute("GET", path, handler)
}
