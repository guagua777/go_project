package channel

func main() {
	// chan int 这是一个整体的类型，可以理解为(chan int)
	// chan <- int 这也是一个整体的类型 (chan <- int)
	// <-chan int 这也是一个整体的类型 (<-chan int)

	// // 无缓冲通道
	// var ch chan int

	// // 有缓冲通道，缓冲区容量为5
	// ch := make(chan int, 5)

	// // 只读通道
	// var recv <-chan int
	// // 只写通道
	// var send chan<- int
}

func f() {
	done := make(chan struct{})
	go func() {
		defer close(done)
		// 业务逻辑
	}()
	// 从done中读取数据
	<-done
}

func g() {
	// struct{} 是空结构体，零内存占用，专门用来做信号通知，不传递实际数据
	done := make(chan struct{})
	go func() {
		done <- struct{}{} // 发一个信号
	}()
	<-done
}

// 第33行为什么要有两个大括号？
// 20:38
// 这是一个很好的 Go 语言细节问题！

// 第 33 行的代码是：

// Go

// done <- struct{}{}
// 这里两对大括号的含义不同：

// 第一对 {} — 属于类型定义部分：struct{} 是 Go 中的空结构体类型（一个没有任何字段的结构体）。
// 第二对 {} — 属于复合字面量（composite literal）部分：用来创建该类型的一个实例。
// 所以 struct{}{} 的完整解读是：创建一个 struct{} 类型的实例。

// 这和普通结构体的写法是同一个模式，只是类型名比较特殊：

// Go

// type Person struct { Name string }

// p := Person{}      // Person 是类型名，{} 是实例化
// s := struct{}{}    // struct{} 是类型名，{} 是实例化
// 为什么用 struct{} 而不是 int 或 bool？
// done 通道的类型是 chan struct{}，用空结构体作为通道元素类型有两个好处：

// 零内存占用：struct{} 不占用任何内存（大小为 0 字节），而 int 占 8 字节，bool 占 1 字节。
// 语义清晰：明确表达这个通道只用来发信号，不传递任何数据。收到 struct{}{} 本身没有意义，有意义的是"收到了"这件事。

type Signal struct{}

func h() {
	// 其中Signal即为struct{}
	s := Signal{}
	println(s)
}
