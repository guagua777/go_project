package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/guagua777/distributed/grades"
	"github.com/guagua777/distributed/log"
	"github.com/guagua777/distributed/registry"
	"github.com/guagua777/distributed/service"
)

func main() {
	host, port := "localhost", "6000"

	serviceAddress := fmt.Sprintf("http://%v:%v", host, port) // localhost:6000

	r := registry.Registration{
		ServiceName: registry.GradingService,
		// 对外提供的服务URL
		ServiceURL: serviceAddress,
		// 初始化一个切片
		RequiredServices: []registry.ServiceName{registry.LogService},
		// 自身url，与注册中心进行沟通
		ServiceUpdateURL: serviceAddress + "/services",
	}

	// 启动服务
	ctx, err := service.Start(
		context.Background(),
		r,
		host,
		port,
		grades.RegisterHandlers,
	)

	// 获取logservice的URL，用于设置客户端的Logger
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil { // 通过GetProvider找到registry.LogService对应的URL
		fmt.Printf("Loggin service found at %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName) // 设置客户端的Logger，其实logProvider就是registry.LogService的URL
	}

	if err != nil {
		// 写日志，该日志会被写到服务其中
		stlog.Fatalln(err)
	}

	// fmt.Printf("local log ......")
	// 写日志，该日志会被写到日志服务中，而不是本地
	stlog.Println("Grading service started.")

	// 在ctx未被关闭之前，ctx.Done()返回的chan struct{}是空的（没有值），所以<-ctx.Done()会阻塞
	// 关闭一个 channel 后，仍然可以从 channel 中读取已经发送的所有元素，直到 channel 中的元素被消耗完为止。
	// 关闭 channel 的行为会使所有的读取操作立即返回零值，并且即使 channel 是空的，也会返回零值。
	<-ctx.Done() // 如果启动http服务器时出现错误/手动停止，service.go中两个协程中就会调用cancel()，就会发送信号

	fmt.Println("Shutting down grading service.")
}
