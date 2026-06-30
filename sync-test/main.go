package main

import (
	"fmt"
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 函数退出自动计数器-1
	fmt.Printf("协程 %d 执行完毕\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)         // 每启动一个协程，计数+1
		go worker(i, &wg) // 必须传指针！
	}

	wg.Wait() // 阻塞，直到所有Done，计数器归0
	fmt.Println("所有协程执行完成")
}
