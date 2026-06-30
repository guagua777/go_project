package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wait = sync.WaitGroup{}

func main() {
	t1 := time.Now()
	wait.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	// 子goroutine中执行业务逻辑
	go func() {
		ip, err := GetIp(ctx)
		fmt.Println(ip, err)
	}()
	wait.Add(1)

	// 子goroutine中执行取消操作
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
		wait.Done()
	}()

	// 主线程
	wait.Wait()
	fmt.Println("执行完成", time.Since(t1))
}

func GetIp(ctx context.Context) (ip string, err error) {

	// 业务逻辑中添加获取取消信号的逻辑
	go func() {
		defer wait.Done()
		select {
		case <-ctx.Done():
			fmt.Println("取消", ctx.Err().Error())
			err = ctx.Err()
			// wait.Done()
			return
		}
	}()

	wait.Add(1)
	defer wait.Done()
	time.Sleep(3 * time.Second)
	ip = "192.168.200.1"
	// wait.Done()
	return
}
