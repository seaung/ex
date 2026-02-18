package main

import (
	"log"
	"time"

	"github.com/seaung/ex"
)

func logger() ex.HandlerFunc {
	return func(ctx *ex.Context) {
		start := time.Now()
		ctx.Next()
		log.Println("time := ", time.Since(start))
	}
}

func main() {
	engine := ex.NewEngine()
	engine.Use(logger())
	engine.GET("/hello", func(ctx *ex.Context) {
		ctx.String(200, "hello ex!!!")
	})
	engine.Run(":9527")
}
