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
