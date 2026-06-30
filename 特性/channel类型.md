// 只读通道
var recv <-chan int


// Done()的返回类型为 只读通道
Done() <-chan struct{}

// 使用
done := Done()
r := <-done



Done() <-chan struct{}
后面的<-chan struct{}是类型吗？
20:54
是的，<-chan struct{} 是一个类型，它是 Done() 方法的返回类型。

我们来拆解这行接口定义：


Go

Done() <-chan struct{}
Done() — 方法名和参数列表（无参数）
<-chan struct{} — 返回类型
<-chan struct{} 这个类型怎么理解？
可以从右往左读：

部分	含义
struct{}	空结构体类型
chan struct{}	传输 struct{} 的通道类型
<-chan struct{}	只读通道类型（只能从中接收，不能发送）
所以 <-chan struct{} 是一个整体类型，表示"元素类型为 struct{} 的只读通道"。

对比三种通道类型

Go

chan struct{}      // 可读可写通道
chan<- struct{}    // 只写通道（只能发送，不能接收）
<-chan struct{}    // 只读通道（只能接收，不能发送）
为什么 Done() 返回只读通道？
这是 Go 的一个接口设计惯例：调用者只能从 Done 通道中接收信号（等待取消），而不能发送信号。取消的权力只属于 Context 的创建者（内部实现），外部不应该能往里发数据。

如果返回 chan struct{}（可读可写），调用者就能意外地 ctx.Done() <- struct{}{}，这会破坏 Context 的语义。返回 <-chan struct{} 在类型层面就禁止了这种误用。