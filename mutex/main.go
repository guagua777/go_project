package main

import (
	"fmt"
	"sync"
)

var (
	counter int
	mu      sync.Mutex
	wg      sync.WaitGroup
)

func add() {
	defer wg.Done()
	mu.Lock()         // 加锁
	defer mu.Unlock() // 函数退出自动解锁
	counter++
}

func main() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go add()
	}
	wg.Wait()
	fmt.Println(counter) // 稳定输出 1000
}
