package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/guagua777/distributed/registry"
)

func main() {
	// 注册handler
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server

	// const ServicePort = ":3000"
	srv.Addr = registry.ServicePort

	go func() {
		// 启动 HTTP 服务器并将其运行日志记录下来。如果 ListenAndServe 返回一个错误，这个错误将被记录到日志中
		// 这个方法是阻塞的，也就是说，它会阻塞当前的 goroutine 直到服务器被显式地关闭或遇到一个致命的错误
		log.Println(srv.ListenAndServe())
		// 除非服务器被显式地关闭或遇到一个致命的错误，否则srv.ListenAndServe()会一直监听
		cancel()
	}()

	go func() {
		fmt.Println("Registry service started. Press any key to stop.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done() // 这行代码会阻塞，直到上下文被取消（也就是 cancel() 被调用）
	fmt.Println("Shutting down registry service")
}
