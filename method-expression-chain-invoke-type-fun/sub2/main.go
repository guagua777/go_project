package main

import (
	"net/http"
)

// 具名函数类型
type HttpHandler func(w http.ResponseWriter, r *http.Request)

func (h HttpHandler) WithLog() HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		println("记录请求日志")
		h(w, r)
	}
}

func (h HttpHandler) WithAuth() HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		println("校验token")
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("hello world"))
}

func main() {
	route := HttpHandler(hello).WithLog().WithAuth()
	http.HandleFunc("/", route)
	_ = http.ListenAndServe(":8080", nil)
}
