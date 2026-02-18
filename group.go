package ex

/**
 * 路由分组
 */
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func (rg *RouterGroup) AddGroup(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		parent: rg,
		engine: engine,
	}
	rg.engine.groups = append(rg.engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}

func (rg *RouterGroup) addRoute(method, path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	finalHandlers := make([]HandlerFunc, 0)
	for g := rg; g != nil; g = g.parent {
		finalHandlers = append(g.middlewares, finalHandlers...)
	}
	finalHandlers = append(finalHandlers, handlers...)
	rg.engine.router.addRoute(method, fullPath, finalHandlers)
}

func (rg *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("GET", fullPath, handlers...)
}

func (rg *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("POST", fullPath, handlers...)
}

func (rg *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("DELETE", fullPath, handlers...)
}

func (rg *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("PUT", fullPath, handlers...)
}

func (rg *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("OPTIONS", fullPath, handlers...)
}

func (rg *RouterGroup) HEAD(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("HEAD", fullPath, handlers...)
}

func (rg *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.addRoute("PATCH", fullPath, handlers...)
}

func (rg *RouterGroup) Any(path string, handlers ...HandlerFunc) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	for _, method := range methods {
		rg.addRoute(method, path, handlers...)
	}
}
