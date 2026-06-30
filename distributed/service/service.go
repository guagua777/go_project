package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

/**
	整体内容：
	1. 注册服务，即注册handler
	2. 启动服务
**/

// 向http中注册handler
// 公共函数，用于启动服务
func Start(ctx context.Context, serviceName, host string, port string,
	registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, serviceName, host, port)

	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server // 定义一个http服务器
	srv.Addr = ":" + port

	// 使用了两个 goroutine 来并发地处理 HTTP 服务器的运行和停止请求
	go func() { // 启动http服务器
		log.Println(srv.ListenAndServe())
		cancel() // 当服务器退出时，取消 context
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx) // 用户输入任意内容后关闭服务器
		cancel()          // 取消context
	}()
	return ctx
}
