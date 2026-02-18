package ex

/*
 * ex框架的路由处理逻辑
 */

// 路由结构体
type Router struct {
	handlers map[string]HandlerFunc
}

// 实例化路由结构体
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

// 用户处理用户请求
func (rt *Router) handle(ctx *Context) {
	key := ctx.Req.Method + "-" + ctx.Req.URL.Path
	if handler, ok := rt.handlers[key]; ok {
		ctx.handlers = append(ctx.handlers, handler)
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(404, "404 NOT FOUND")
		})
	}
	ctx.index = -1
	ctx.Next()
}

// 用于注册用户路由操作
func (rt *Router) addRoute(method, path string, handler HandlerFunc) {
	key := method + "-" + path
	rt.handlers[key] = handler
}
