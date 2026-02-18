package main

import (
	"github.com/seaung/ex"
)

func main() {
	engine := ex.NewEngine()
	engine.GET("/hello", func(ctx *ex.Context) {
		ctx.String(200, "hello world!!!")
	})
	engine.Run(":9527")
}
