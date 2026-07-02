package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	go func() {
		fmt.Println("服务器启动在 :8080")
		http.ListenAndServe(":8080", nil)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/hello")
	fmt.Printf("注册 /hello 前: 状态码=%v, 错误=%v\n", resp.StatusCode, err)
	if resp != nil {
		resp.Body.Close()
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello!")
	})
	fmt.Println("已注册 /hello")

	time.Sleep(100 * time.Millisecond)

	resp, err = http.Get("http://localhost:8080/hello")
	fmt.Printf("注册 /hello 后: 状态码=%v, 错误=%v\n", resp.StatusCode, err)
	if resp != nil {
		resp.Body.Close()
	}
}
