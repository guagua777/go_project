package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker1(ctx context.Context) {
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

func main1() {
	// 创建可取消上下文
	ctx, cancel := context.WithCancel(context.Background())

	go worker1(ctx)

	// 3秒后手动取消
	time.Sleep(3 * time.Second)
	cancel()

	// 等待观察goroutine退出
	time.Sleep(1 * time.Second)
	fmt.Println("主程序结束")
}

func worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
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

// case <-ctx.Done():
// 这个只读的返回值，不需要一个变量来接收吗
// 13:28
// 这是一个很好的问题！

// ctx.Done() 的返回值
// ctx.Done() 返回的类型是 <-chan struct{}，即一个只读的空结构体通道。

// 关键点在于：struct{} 是空结构体，它不携带任何数据，大小为 0 字节。这个 channel 的作用纯粹是信号通知——能从里面读出值，就说明 context 被取消了或超时了，值本身毫无意义。

// 所以有两种写法：

// Go

// // 写法1：丢弃返回值（最常见，推荐）
// case <-ctx.Done():
//     fmt.Println("收到取消信号")

// // 写法2：用变量接收，但 v 是 struct{}{}，毫无用处
// case v := <-ctx.Done():
//     fmt.Println("收到取消信号", v)  // v = struct{}{}，没有任何信息
// 类比理解
// 可以把它类比为日常生活中的信号灯：

// 🟢 绿灯亮了 → 你只需要知道"可以走了"，不需要灯给你什么数据
// ctx.Done() 关闭 → 你只需要知道"该停了"，不需要它告诉你为什么停
// 如果想知道取消的原因，应该用 ctx.Err()：

// Go

// case <-ctx.Done():
//     fmt.Println("取消原因:", ctx.Err())  // context canceled 或 context deadline exceeded
// 总结
// 方式	用途
// <-ctx.Done()	只关心"是否取消"的信号
// ctx.Err()	关心取消的具体原因
// v := <-ctx.Done()	没必要，v 永远是 struct{}{}
// 所以 case <-ctx.Done(): 不用变量接收，是因为这个 channel 本身就是设计为纯信号通道，返回值无任何信息量，丢弃即可。

func worker2(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("任务退出，原因：", ctx.Err())
			return
		default:
			fmt.Println("正在工作2...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// 使用sync wait修改
	// 创建可取消上下文
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	wg.Add(1)
	go worker(ctx, &wg)

	wg.Add(1)
	go worker2(ctx, &wg)

	// 3秒后手动取消
	time.Sleep(3 * time.Second)
	// 该context下的所有的任务都会取消
	cancel()

	// 等待观察goroutine退出
	// time.Sleep(1 * time.Second)
	wg.Wait()
	fmt.Println("主程序结束")
}
