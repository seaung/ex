package main

import (
	"net/http"

	"github.com/seaung/ex"
)

func main() {
	engine := ex.NewEngine()
	engine.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	engine.Run(":9527")
}
