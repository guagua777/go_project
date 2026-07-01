package main

import (
	"runtime"
	"sync"
)

// ReentrantMutex 可重入互斥锁
type ReentrantMutex struct {
	mu      sync.Mutex
	ownerID int64 // 持有锁的goroutine id
	recur   int   // 重入次数
}

// Lock 加锁：支持同一协程多次调用
func (r *ReentrantMutex) Lock() {
	gid := goroutineID()

	// 当前协程已经持有锁，只增加重入计数
	if r.ownerID == gid {
		r.recur++
		return
	}

	// 抢占底层锁
	r.mu.Lock()
	r.ownerID = gid
	r.recur = 1
}

// Unlock 解锁：必须和Lock成对调用
func (r *ReentrantMutex) Unlock() {
	gid := goroutineID()
	if r.ownerID != gid {
		panic("sync: unlock of mutex not held by goroutine")
	}

	r.recur--
	if r.recur == 0 {
		// 重入次数归零，真正释放锁
		r.ownerID = 0
		r.mu.Unlock()
	}
}

// goroutineID 获取当前goroutine唯一ID
func goroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	println(string(buf[:n]))
	var id int64
	for i := 0; i < n; i++ {
		if buf[i] == 'g' && i+1 < n {
			// 解析g后面的数字
			for i++; i < n && buf[i] >= '0' && buf[i] <= '9'; i++ {
				id = id*10 + int64(buf[i]-'0')
			}
			break
		}
	}
	return id
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		id := goroutineID()
		println(id)
	}()

	wg.Wait()
}
