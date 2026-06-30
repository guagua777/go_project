context.WithCancel解释下
context.WithCancel 详解
一、基础定义
go
运行
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
WithCancel 接收一个父上下文，返回两个值：
子上下文 ctx：携带截止信号
取消函数 cancel()：主动关闭上下文
作用：手动控制协程退出，用于 goroutine 的并发控制、任务终止。
二、核心特性
调用 cancel() 后
ctx.Done() 通道立刻关闭
ctx.Err() 返回 context.Canceled
父上下文取消时，子上下文会跟着自动取消（上下文链式传递）
cancel 函数必须调用，哪怕任务正常完成，避免内存泄漏。
多次调用 cancel() 不会 panic，只是无效执行。
三、最简示例
go
运行
package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("任务退出，原因：", ctx.Err())
			return
		default:
			fmt.Println("正在工作...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// 创建可取消上下文
	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx)

	// 3秒后手动取消
	time.Sleep(3 * time.Second)
	cancel()

	// 等待观察goroutine退出
	time.Sleep(1 * time.Second)
	fmt.Println("主程序结束")
}
运行结果：
plaintext
正在工作...
正在工作...
正在工作...
任务退出，原因：context canceled
主程序结束
四、Done 与 Err 机制
ctx.Done()：返回一个只读 channel，上下文未取消时阻塞，取消后通道关闭。
ctx.Err()：
未取消：返回 nil
手动调用 cancel：context.Canceled
超时取消：context.DeadlineExceeded
五、父子上下文链式取消
go
运行
parent, parentCancel := context.WithCancel(context.Background())
child, _ := context.WithCancel(parent)

parentCancel()
// child会被自动取消
只要祖先上下文被取消，所有后代上下文全部自动取消，非常适合服务级联关闭。
六、最佳实践
函数收到 context.Context 必须优先监听 ctx.Done()
只要调用了 WithCancel，就一定要 defer cancel()
go
运行
func doTask() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 防止泄漏
	// ...业务逻辑
}
不要把 nil Context 传入函数，默认用 context.Background() 作为根上下文。
七、适用场景
主动终止 goroutine
HTTP 请求中途断开，关闭下游调用
关闭子协程、清理资源
服务退出时批量终止所有业务协程
需要我再对比一下 WithCancel / WithTimeout / WithDeadline 的区别吗？







1. 三者签名对比
go
运行
// 手动取消
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

// 相对超时（最常用）
func WithTimeout(parent Context, timeout time.Duration) (ctx Context, cancel CancelFunc)

// 绝对截止时间
func WithDeadline(parent Context, d time.Time) (ctx Context, cancel CancelFunc)
2. 核心区别
① WithCancel
触发时机：手动调用 cancel () 才会取消
错误原因：ctx.Err() = context.Canceled
适用：不确定等待时长，由外部信号终止协程（用户关闭连接、收到退出信号）
② WithTimeout
本质：对 WithDeadline 的封装
触发时机：要么时间到自动取消，要么手动 cancel
超时后错误：context.DeadlineExceeded
适用：RPC、数据库请求、HTTP 调用，防止无限阻塞
③ WithDeadline
触发时机：到达某个固定时间点自动取消
适用：多个任务共用同一个截止时刻（分布式调用统一超时）
3. 关键共性（三条都遵守）
父上下文一旦取消，子上下文全部自动取消。
都必须调用 cancel，建议 defer cancel()，防止 goroutine 泄漏。
Done () 通道关闭后，Err () 才会返回非 nil 错误。
多次调用 cancel 是安全的，不会 panic。
4. 代码示例对照
WithCancel
go
运行
ctx, cancel := context.WithCancel(context.Background())
time.Sleep(2 * time.Second)
cancel() // 手动终止
WithTimeout
go
运行
// 3秒后自动超时
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
WithDeadline
go
运行
// 固定时刻到期
deadline := time.Now().Add(3 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
5. 错误区分
go
运行
select {
case <-ctx.Done():
    err := ctx.Err()
    if err == context.Canceled {
        fmt.Println("主动取消")
    } else if err == context.DeadlineExceeded {
        fmt.Println("超时")
    }
}
6. 场景选型总结
表格
函数	何时用
WithCancel	手动控制退出，无固定超时
WithTimeout	绝大多数网络调用，设置等待时长
WithDeadline	全链路统一截止时间
7. 易错点
父上下文超时，子上下文也会被连带取消。
如果不希望被子父上下文影响截止时间，只能新建根上下文 context.Background()。
要不要我再讲一下 context 传参规范，以及不要把 value 上下文滥用的坑？