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
	ctx := context.Background()
	wait.Add(1)
	go func() {
		defer wait.Done()
		ip, err := GetIp(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ip)
	}()
	wait.Wait()
	fmt.Println("执行完成", time.Since(t1))
}

func GetIp(ctx context.Context) (ip string, err error) {
	time.Sleep(2 * time.Second)
	ip = "192.168.200.1"
	return
}
