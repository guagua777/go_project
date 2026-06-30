go语言的sync.WaitGroup{}，介绍下
sync.WaitGroup 完整介绍
一、作用
sync.WaitGroup 用来等待一组 Goroutine 全部执行完毕，主线程阻塞等待所有子协程完成，再继续向下执行。
本质是一个计数器。
二、三个核心方法
go
运行
var wg sync.WaitGroup
Add(delta int)
增加 / 减少等待计数器的值。
开启协程前：wg.Add(1)，计数器 + 1
禁止在协程内部调用 Add，容易出现竞态。
Done()
等价于 wg.Add(-1)，协程执行完毕时调用，计数器减 1。
一般放在协程函数末尾，搭配 defer wg.Done()。
Wait()
阻塞当前 goroutine，直到计数器变为 0。
三、基础示例
go
运行
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
		wg.Add(1) // 每启动一个协程，计数+1
		go worker(i, &wg) // 必须传指针！
	}

	wg.Wait() // 阻塞，直到所有Done，计数器归0
	fmt.Println("所有协程执行完成")
}
关键点：必须传递 *sync.WaitGroup 指针，传值会拷贝结构体，Done 作用在副本上，永久阻塞。
四、常见坑（重点）
1. 不传指针，只传值
WaitGroup 是值类型，函数传参拷贝副本，子协程修改副本，主协程的计数器永远不为 0，死锁。
2. Add 在协程内部调用
go
运行
// 错误写法
go func() {
	wg.Add(1)
	defer wg.Done()
}()
主协程可能先执行到 wg.Wait()，此时计数器还是 0，Wait 直接放行，协程还没来得及执行。
✅ 规范：Add 必须在 go 语句之前执行。
3. 计数器变成负数，触发 panic
多次调用 Done 或者 Add 负数过多：
go
运行
wg.Add(1)
wg.Done()
wg.Done() // 计数器 < 0，直接 panic
4. Wait 复用问题
WaitGroup 计数器归零之后可以再次复用，但不能在 Wait 还没结束时调用 Add。
五、内部简单原理
sync.WaitGroup 底层结构体包含一个 state 复合变量：
高 32 位：等待者数量（正在阻塞 Wait 的 goroutine）
低 32 位：任务计数器
通过原子操作修改计数，无锁高性能。
六、使用规范总结
启动 goroutine 前调用 wg.Add(n)；
协程内部使用 defer wg.Done()；
传参永远传递 *sync.WaitGroup 指针；
保证计数器最终归零，避免负数 panic；
Wait() 等待所有任务完成。
需要我给你对比一下 WaitGroup、channel + close、errgroup 三者的区别吗？




WaitGroup / 关闭 Channel /errgroup 对比
1. 三者核心用途
表格
方案	核心能力	适用场景	能否传递错误
sync.WaitGroup	只等待协程全部完成	单纯等待，不需要收集错误	❌ 不能传递错误
channel + close	等待 + 结果收集	收集协程返回数据，全部跑完再关闭	可以自己封装错误
golang.org/x/sync/errgroup	等待 + 自动捕获第一个错误	并发任务，一旦出错立刻取消剩余协程	✅ 原生支持 error
2. 代码示例
（1）sync.WaitGroup
只做等待，无法传递错误。
go
运行
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}
	wg.Wait()
}
缺点：
协程内部发生错误，主线程感知不到；
无法中途取消其他协程。
（2）Channel + close 方案
适合收集任务结果。
go
运行
func main() {
	ch := make(chan error, 3)
	for i := 0; i < 3; i++ {
		go func() {
			ch <- nil // 成功写入nil，失败写入err
		}()
	}

	// 等待全部写完再关闭
	for i := 0; i < 3; i++ {
		<-ch
	}
	close(ch)
}
优点：可以接收每个协程的返回值；
缺点：需要手动计数、手动处理取消，代码繁琐。
（3）errgroup（推荐业务并发使用）
errgroup 内部封装了 WaitGroup + context。
只要任意一个协程返回 error，上下文立刻取消，其余协程收到 ctx.Done 信号。
go
运行
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < 3; i++ {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// 业务逻辑
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("任务失败:", err)
	}
}
关键特性：
只要有一个任务返回 error，g.Wait() 立即返回该错误；
自动关闭 context，其他协程可以监听 ctx 实现快速退出；
不用手动写 Add/Done，不容易出错。
3. 核心区别总结
1）WaitGroup
纯计数器，无上下文、无错误传递；
只能等所有协程跑完，不能提前终止；
适合简单的批量并发，无异常处理场景。
2）Channel
可以接收每个协程的数据 / 错误；
没有内置取消机制，要自己实现退出逻辑；
适合需要收集多条返回结果的场景。
3）errgroup（业务首选）
WaitGroup + context 组合封装；
一旦出现错误，自动触发上下文取消，快速熔断；
省去手动 Add/Done，代码简洁，不容易死锁；
绝大多数接口并发、批量任务优先用 errgroup。
4. 选型建议
简单批量并发，不需要错误处理：用 sync.WaitGroup
需要收集每个协程的执行结果：用 buffered channel
并发任务，遇到错误要立刻终止其他协程：用 errgroup.Group
要不要我再补充 errgroup 源码精简解读，看它是怎么结合 WaitGroup + Context 实现熔断的？