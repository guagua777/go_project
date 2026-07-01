sync.Mutex（Go 互斥锁）
一、基础说明
包路径：sync
结构体：sync.Mutex，排他锁，同一时间只允许一个 Goroutine 持有锁。
核心方法
Lock()
获取锁；锁已被占用时，当前 Goroutine 阻塞等待。
Unlock()
释放锁。
注意：
必须和 Lock() 成对调用；
禁止在未加锁时解锁，会 panic；
不要在不同 Goroutine 里加锁、解锁。
最简示例
go
运行
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
二、底层结构（src/sync/mutex.go）
go
运行
type Mutex struct {
	state int32  // 状态位：锁占用、唤醒、等待goroutine计数
	sema  uint32 // 信号量，用于Goroutine阻塞休眠
}
state 字段拆分：
bit0：是否已上锁（locked）
bit1：是否有 Goroutine 被唤醒（woken）
高位：等待队列 Goroutine 数量
锁抢占逻辑（饥饿模式 / 正常模式）
正常模式（非饥饿）
新来的 Goroutine 可以直接自旋抢锁，排在等待队列前面，性能高，可能导致老等待者长期抢不到锁（饥饿）。
饥饿模式
一旦等待超过 1ms，Mutex 进入饥饿模式；
新来的 Goroutine 不许自旋，直接排到等待队尾，严格 FIFO，保证公平。
退出饥饿：当一个 Goroutine 抢到锁后，队列后面已经没人等待，切回正常模式。
三、常见坑
忘记 Unlock
死锁，后续所有 Lock 永久阻塞。
最佳写法：mu.Lock(); defer mu.Unlock()
重复 Lock（同一协程）
go
运行
mu.Lock()
mu.Lock() // 死锁！Mutex不可重入
Mutex 不支持递归重入；需要重入锁请用 sync.RWMutex 也不行，必须自己实现可重入锁（记录持有者 GoroutineID）。
跨协程解锁
go
运行
go func(){ mu.Lock() }()
mu.Unlock() // panic: sync: unlock of unlocked mutex
Unlock 放在 defer 之外，提前退出漏掉释放锁
四、Mutex vs RWMutex
表格
类型	读	写	场景
Mutex	互斥	互斥	读写均衡、简单竞争
RWMutex	共享并发读	排他写	读多写少
五、面试高频问题
Mutex 是否可重入？
不可重入，同一个 Goroutine 连续两次 Lock 会死锁。
自旋什么时候开启？
正常模式下，锁被占用时，会短暂自旋几次，不立刻陷入休眠，减少上下文切换。
饥饿模式解决了什么问题？
防止新来协程不断抢占锁，导致队列尾部的 Goroutine 永远得不到执行。
为什么 Unlock 可以被其他协程执行？语法允许，但业务严禁这么做。
六、可重入 Mutex 简易实现（拓展）
利用 goroutine id 记录持有者：
go
运行
type ReentrantMutex struct {
	mu    sync.Mutex
	owner int64
	count int
}
需要我把完整可重入锁代码给你吗？




跨协程解锁
go func(){ mu.Lock() }()
mu.Unlock() // panic: sync: unlock of unlocked mutex

这个的问题是什么
1. 直接报错原因
go
运行
go func(){ mu.Lock() }()
mu.Unlock()

主线程执行 mu.Unlock() 的时候，锁根本还没有被主线程持有。
Go 的 sync.Mutex 有一条强制规则：
只能由当前持有锁的 Goroutine 来调用 Unlock ()
对未上锁的锁执行 Unlock，直接触发 panic：
sync: unlock of unlocked mutex
2. 时序问题（更深一层）
执行顺序大概率是这样：
启动子协程，进入就绪队列
主线程继续执行，立刻调用 mu.Unlock()
此时子协程还没调度执行到 mu.Lock()，锁处于未锁定状态
→ 解锁空锁，panic。
就算调度变慢，子协程先抢到 CPU 执行了 mu.Lock()：
锁是子协程持有
主线程执行 mu.Unlock()
依然违反规则：非持有者不能释放锁，照样 panic。
3. Mutex 内部是怎么校验的？
Mutex.state 标记锁状态，但不会记录持有锁的 GoroutineID。
它只做两件检查：
如果锁处于未锁定状态：Unlock 直接 panic；
如果锁已经锁定：仅仅修改状态释放，不会校验调用者是不是持有者。
所以：
场景 A：锁还没 Lock → Unlock → panic（你例子里的情况）
场景 B：别的协程加了锁，你在外部 Unlock
go
运行
go func(){mu.Lock(); time.Sleep(time.Second)}()
time.Sleep(10 * time.Millisecond)
mu.Unlock() // 不会panic，但属于非法用法，破坏并发安全

结论：语法层面允许跨协程解锁，但属于非法编程模式，会破坏锁的所有权，造成数据竞争。官方强烈禁止。
4. 总结两点核心问题
时序竞争：主线程先执行 Unlock，此时锁尚未被锁住，触发 “解锁未上锁的互斥锁” panic；
锁所有权违规：即使时序调换，锁由子协程持有，主线程去释放锁，违背了 “谁加锁，谁解锁” 的原则，会造成不可预知的并发错乱。
需要我再给你一段演示 “跨协程解锁无 panic 但引发脏数据” 的代码吗？



可重入互斥锁完整实现
原理：
底层包裹普通 sync.Mutex
记录当前持有锁的 goroutine ID
记录重入次数，同一协程多次加锁只增加计数，不会阻塞
必须层层解锁，计数归零才真正释放底层锁
go
运行
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
测试用例
go
运行
func main() {
	var lock ReentrantMutex

	lock.Lock()
	lock.Lock() // 第二次重入，不会死锁
	println("重入次数：", lock.recur)

	lock.Unlock()
	println("解锁一次，剩余次数：", lock.recur)

	lock.Unlock() // 计数归零，释放底层锁
}
关键点说明
原生 sync.Mutex 没有保存 goroutine ID，所以天然不可重入；
依靠解析栈文本拿到 goroutine ID，这是标准实现方案；
严格遵守：谁加锁，谁解锁，依然不能跨协程释放；
不能无限调用 Lock，必须保证 Lock 和 Unlock 数量对等，否则锁永远不会释放。
缺陷提醒
goroutineID 通过解析字符串获取，属于非官方私有实现，Go 版本迭代基本稳定，生产环境可以正常使用。
要不要我再写一版加上 TryLock 非阻塞尝试加锁？