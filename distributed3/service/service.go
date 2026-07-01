package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/guagua777/distributed/registry"
)

/**
	整体内容：
	1. 注册服务，即注册handler
	2. 启动服务
**/

// 向http中注册handler
// 公共函数，用于启动服务
/**
由三个部分组成：
1. 对应服务的handler
2. 启动相应的服务（根据host和端口）
3. 注册服务（将服务注册信息发送给注册服务）

对应的日志即为：
1. 日志服务的handler
2. 启动日志服务（根据host和端口）
3. 注册日志服务（将日志服务的注册信息发送给注册服务）

**/
func Start(ctx context.Context, reg registry.Registration, host string, port string,
	registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	// 启动服务后，注册服务
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server // 定义一个http服务器
	srv.Addr = ":" + port

	// 使用了两个 goroutine 来并发地处理 HTTP 服务器的运行和停止请求
	go func() { // 启动http服务器
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel() // 当服务器退出时，取消 context
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx) // 用户输入任意内容后关闭服务器
		cancel()          // 取消context
	}()
	return ctx
}

// 是的，你的理解是对的！srv.ListenAndServe() 会阻塞，但 cancel() 最终还是会执行。让我详细解释：

// ListenAndServe 的行为
// ListenAndServe 是一个阻塞调用，它会一直运行直到：

// 服务器被 Shutdown 主动关闭
// 发生致命错误（如端口被占用）
// 所以正常情况下，cancel() 确实不会立刻执行，而是等到 ListenAndServe 返回后才执行。

// 但为什么这样写没问题？
// 关键在于这段代码在 goroutine 里，而且有另一个 goroutine 负责关闭服务器：

// go
// // goroutine 1: 启动服务器（阻塞在这里）
// go func() {
//     log.Println(srv.ListenAndServe()) // 阻塞，直到服务器关闭
//     cancel() // 服务器关闭后才执行
// }()

// // goroutine 2: 等待用户输入来关闭服务器
// go func() {
//     fmt.Printf("%v started. Press any key to stop. \n", serviceName)
//     var s string
//     fmt.Scanln(&s)          // 阻塞，等待用户输入
//     srv.Shutdown(ctx)       // 用户输入后，关闭服务器 → 让 goroutine 1 的 ListenAndServe 返回
//     cancel()
// }()
// 执行流程

// plainText
// 时间线：
// ─────────────────────────────────────────────────────►

// goroutine 1:  ListenAndServe(阻塞中) ──────────────→ 返回 → cancel()
//                                               ↑
// goroutine 2:  Scanln(等待输入) ──→ Shutdown ──┘
// 两个 goroutine 同时启动，各自阻塞
// 用户按任意键 → Scanln 返回 → 调用 srv.Shutdown(ctx)
// Shutdown 会让 ListenAndServe 停止监听并返回 → cancel() 被执行
// 小问题
// 这里其实有一个潜在问题：srv.Shutdown(ctx) 传入的 ctx 可能已经被 cancel() 取消了。因为两个 goroutine 都调用了 cancel()，如果 goroutine 1 的 ListenAndServe 因其他原因先返回并调用了 cancel()，那么 goroutine 2 中的 srv.Shutdown(ctx) 使用的就是一个已取消的 context，这可能导致优雅关闭失败。

// 更健壮的写法可以用 context.Background() 或独立的 context 来调用 Shutdown，不过对于这个简单的示例来说，当前写法在正常流程下是可以工作的。
