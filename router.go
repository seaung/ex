package ex

/*
 * ex框架的路由处理逻辑
 */

import "net/http"

// 路由结构体
type Router struct {
	handlers map[string]http.HandlerFunc
}

// 实例化路由结构体
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]http.HandlerFunc),
	}
}

// 用户处理用户请求
func (rt *Router) handle(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := rt.handlers[key]; ok {
		handler(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Server Not Found"))
}

// 用于注册用户路由操作
func (rt *Router) addRoute(method, path string, handler http.HandlerFunc) {
	key := method + "-" + path
	rt.handlers[key] = handler
}
