package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/guagua777/distributed/grades"
	"github.com/guagua777/distributed/registry"
	"github.com/guagua777/distributed/service"
)

func main() {
	host, port := "localhost", "6000"

	serviceAddress := fmt.Sprintf("http://%v:%v", host, port) // localhost:6000

	r := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddress,
	}

	// 启动服务
	ctx, err := service.Start(
		context.Background(),
		r,
		host,
		port,
		grades.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	// 在ctx未被关闭之前，ctx.Done()返回的chan struct{}是空的（没有值），所以<-ctx.Done()会阻塞
	// 关闭一个 channel 后，仍然可以从 channel 中读取已经发送的所有元素，直到 channel 中的元素被消耗完为止。
	// 关闭 channel 的行为会使所有的读取操作立即返回零值，并且即使 channel 是空的，也会返回零值。
	<-ctx.Done() // 如果启动http服务器时出现错误/手动停止，service.go中两个协程中就会调用cancel()，就会发送信号

	fmt.Println("Shutting down grading service.")
}
