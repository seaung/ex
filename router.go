package ex

/*
 * ex框架的路由处理逻辑
 */

// 路由结构体
type Router struct {
	handlers map[string]map[string]HandlerFunc
}

// 实例化路由结构体
func newRouter() *Router {
	return &Router{
		handlers: make(map[string]map[string]HandlerFunc),
	}
}

// 用户处理用户请求
func (rt *Router) handle(ctx *Context) {
	if methodMap, ok := rt.handlers[ctx.Req.Method]; ok {
		if handler, ok := methodMap[ctx.Req.URL.Path]; ok {
			handler(ctx)
			return
		}
	}
	ctx.String(404, "404 NOT FOUND")
}

// 用于注册用户路由操作
func (rt *Router) addRoute(method, path string, handlers []HandlerFunc) {
	if rt.handlers[method] == nil {
		rt.handlers[method] = make(map[string]HandlerFunc)
	}
	rt.handlers[method][path] = func(ctx *Context) {
		ctx.handlers = handlers
		ctx.index = -1
		ctx.Next()
	}
}
