线程间通讯，通过channel来通讯，而不是通过共享内存来通讯。


一、基础特性（对比 OS 线程）
内存极小
Goroutine：初始栈只有 2KB，运行时自动伸缩，最大 1GB。
OS 线程：默认栈 1MB～8MB。
一台机器可以轻松创建上百万个 Goroutine，线程只能几千个。
调度在用户态
由 Go 运行时 (runtime) 自行调度，不需要陷入操作系统内核，切换开销远小于线程。
M:N 调度模型
将成千上万个 G（协程）映射到少量 M（操作系统线程），由 P（逻辑处理器）做调度，这就是 GMP 模型。
创建极其简单
go
运行
package main
import "fmt"

func hello() {
    fmt.Println("goroutine running")
}

func main() {
    go hello() // 开启一个协程
    fmt.Scanln() // 防止main退出
}
main本身也是一个 Goroutine。
二、GMP 三大核心组件（调度核心）
1. G = Goroutine
协程本身，保存：程序计数器 PC、栈、寄存器、运行状态（就绪、运行、阻塞、死亡）。
队列分为：
P 的本地队列（优先执行）
全局队列
2. M = Machine
对应操作系统内核线程，是真正执行代码的载体。
M 必须绑定 P 才能运行 G；如果 G 发起系统调用阻塞，M 被卡住，P 会立刻解绑并新建 M 承接剩余任务，保证 CPU 不空闲。
3. P = Processor（逻辑处理器）
调度上下文，是 GMP 的核心创新：
持有运行 Go 代码需要的资源
维护本地协程运行队列（最多 256 个 G）
P 的数量默认等于 CPU 核心数，由环境变量GOMAXPROCS控制。
三者关系
1 个 P 绑定 1 个 M，M 执行 P 队列里的 G；
P 队列为空时会执行工作窃取 (Work Stealing)，从其他 P 偷取 G，充分利用多核 CPU。
三、Goroutine 的几种状态
Runnable 就绪：在队列里等待被调度
Running 运行中：正在 M 上执行
Blocked 阻塞：
锁、channel 等待
网络 IO、文件 IO 系统调用
阻塞时 P 会被释放，去执行其他 G。
Gdead 终止：执行完毕，等待回收栈内存。
四、两大调度机制
1. 工作窃取 Work Stealing
空闲的 P 从其他繁忙 P 的本地队列尾部偷取一半 G，避免 CPU 空转，保证多核负载均衡。
2. 抢占式调度（Go 1.14+）
早期是协作式，死循环协程会霸占 CPU；
现在基于信号实现抢占：单个 G 连续运行超过 10ms，调度器强制剥夺 P，让出 CPU 给其他协程。
3. 阻塞处理（Hand-off 移交）
当 G 执行 syscall（读写文件、网络阻塞），M 阻塞：
P 脱离当前 M
runtime 新建一个 M 绑定 P 继续执行其他 G
IO 完成后，原 G 重新进入就绪队列，保证多核利用率。
五、协程之间通信：Channel
协程不要用共享内存加锁，推荐通信来共享内存：
go
运行
ch := make(chan int)
go func(){ ch <- 100 }() // 子协程写入
val := <-ch // 主协程读取
六、协程优缺点
优点
极高并发，百万协程无压力
用户态切换，性能远超线程
runtime 自动管理栈，自动扩容缩容
配合 channel+select，轻松编写高并发网络程序
限制
Goroutine 只在 Go 运行时内调度，无法被操作系统直接调度
CGO 调用 C 阻塞时，会占用 M 导致 P 被卡住，无法自动让出 CPU
七、常见面试要点精简版
初始栈 2KB，动态扩容缩容；
调度模型：GMP，M:N 多路复用；
P 数量 = GOMAXPROCS (CPU 核数)；
IO 阻塞会新建 M，P 不被阻塞；
队列为空执行工作窃取；
1.14 之后实现信号抢占，杜绝单个 G 长期霸占 CPU；
协程通信优先 channel，少用锁。