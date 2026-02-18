package ex

import (
	"log"
	"net/http"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		log.Printf("[%d] - %s in %v\n", ctx.StatusCode, ctx.Path, time.Since(start))
	}
}

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered: ", err)
				ctx.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}
