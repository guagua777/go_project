package main

import "fmt"

func hello() {
	fmt.Println("goroutine running")
}

func main1() {
	go hello()   // 开启一个协程
	fmt.Scanln() // 防止main退出
}

func main() {
	ch := make(chan int)
	go func() { ch <- 100 }() // 子协程写入
	val := <-ch               // 主协程读取
	fmt.Println(val)
}
