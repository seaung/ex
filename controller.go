package ex

/*
 * 这个是为了让某个具体自定义类型实现Controller接口
 * 然后通过RegisterController方法来注册，根据http方法自动调用具体类型下的方法
 */

import (
	"reflect"
	"strings"
)

type Controller interface {
	Get(ctx *Context)
	Post(ctx *Context)
	Put(ctx *Context)
	Patch(ctx *Context)
	Delete(ctx *Context)
	Options(ctx *Context)
	Head(ctx *Context)
}

/*
- 这个方法用于注册某个实现了Controller接口的具体方法

- 使用示例:

	package main

	import "github.com/seaung/ex"

	type UserController struct{}

	func (u *UserController) Get(ctx *ex.Context) {
	ctx.Json(200, map[string]string{"method": "GET", "msg": "fetch users"})
	}

	func (u *UserController) Post(ctx *ex.Context) {
	var data map[string]string
	if err := ctx.ShouldBindJson(&data); err != nil {
	ctx.Json(400, map[string]string{"error": err.Error()})
	return
	}
	ctx.Json(200, map[string]interface{}{
	"method": "POST",
	"data":   data,
	})
	}

	func (u *UserController) Put(ctx *ex.Context) {
	ctx.Json(200, map[string]string{"method": "GET", "msg": "fetch users"})
	}

	func (u *UserController) Delete(ctx *ex.Context) {
	ctx.Json(200, map[string]string{"method": "GET", "msg": "fetch users"})
	}

	func (u *UserController) Patch(ctx *ex.Context) {
	ctx.Json(200, map[string]string{"method": "GET", "msg": "fetch users"})
	}

	func (u *UserController) Options(ctx *ex.Context) {
	ctx.Json(200, map[string]string{"method": "GET", "msg": "fetch users"})
	}

	func (u *UserController) Head(ctx *ex.Context) {
	ctx.Json(200, map[string]string{"method": "GET", "msg": "fetch users"})
	}

	func main() {

	r := ex.DefaultEngine()
	userCtrl := &UserController{}

	r.AddGroup("/user").RegisterController("", userCtrl)
	r.Run(":9527")
	}
*/
func (rg *RouterGroup) RegisterController(path string, ctrl Controller) {
	if ctrl == nil {
		return
	}

	methods := []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	val := reflect.ValueOf(ctrl)
	for _, method := range methods {
		_func := val.MethodByName(strings.Title(strings.ToLower(method)))
		if _func.IsValid() {
			rg.addRoute(method, path, func(ctx *Context) {
				_func.Call([]reflect.Value{reflect.ValueOf(ctx)})
			})
		}
	}
}
