package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	go GetIp(ctx)

	// 5秒到了，手动结束协程
	time.Sleep(5 * time.Second)
	//cancel() // 可以手动取消，也可让他自然超时

	// 模拟主线程阻塞
	time.Sleep(1 * time.Second)

}

func GetIp(ctx context.Context) {
	fmt.Println("获取ip中")
	// 等待请求完成或者被取消
	select {
	case <-ctx.Done():
		// 请求被取消
		fmt.Println("请求超时或被取消", ctx.Err()) // 可以通过err判断是超时还是取消
	}
}
