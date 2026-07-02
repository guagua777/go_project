// log/server.go中有服务器端的逻辑，但是客户端的服务想使用这个service还是很麻烦
// 为了让客户端的服务方便的使用log/server.go
package log

import (
	"bytes"
	"fmt"
	"io"
	stlog "log"
	"net/http"

	"github.com/guagua777/distributed/registry"
)

// 客户端服务本地写日志的设置
func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)                               // 服务端log设置了时间戳，但客户端不需要
	stlog.SetOutput(&clientLogger{url: serviceURL}) // 客户端记录日志的输出应该指向服务端的logger
}

type clientLogger struct {
	url string
}

var _ io.Writer = (*clientLogger)(nil)

// 这个日志是写到哪里去了
// clientLogger需要实现io.Writer这个接口，所以实现以下Write方法
func (cl clientLogger) Write(data []byte) (int, error) {
	fmt.Printf("Write: %s\n", data)
	b := bytes.NewBuffer([]byte(data))
	// b := bytes.NewBuffer([]byte(p[:len(data) - 1]))
	res, err := http.Post(cl.url+"/log", "text/plain", b) // 写到服务端
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message. Service responed with code %v", res.StatusCode)
	}
	return len(data), nil
}

// Chat
// Thought process
// The user is asking about the Go idiom var _ io.Writer = (*clientLogger)(nil). This is a compile-time interface compliance check.

// 这是 Go 中一个非常经典的惯用写法：编译期接口合规性检查。

// 逐部分拆解

// go
// var _ io.Writer = (*clientLogger)(nil)
// 部分	含义
// var _	声明一个变量，名字是 _（空白标识符），意味着这个变量不会被使用，编译器不会报"未使用变量"的错误
// io.Writer	变量的类型是 io.Writer 接口
// =	赋值
// (*clientLogger)(nil)	把 nil 转换成 *clientLogger 类型，得到一个类型为 *clientLogger 的 nil 指针
// 整体	将 *clientLogger 类型的值赋给 io.Writer 类型的变量
// 核心作用
// 如果 *clientLogger 没有实现 io.Writer 接口，这行代码会编译报错。

// 相当于一个"守卫"，确保你的类型确实满足某个接口。比如你删掉了 Write 方法：

// go
// // 如果删掉 Write 方法，编译时会报错：
// //   cannot use (*clientLogger)(nil) (type *clientLogger) as type io.Writer in assignment:
// //       *clientLogger does not implement io.Writer (missing Write method)
// 为什么用 (*clientLogger)(nil) 而不是 &clientLogger{}？
// (*clientLogger)(nil) 只是一个类型转换，不会分配内存，零开销
// &clientLogger{} 会真正创建一个结构体实例，有内存分配开销（虽然很小）
// 等价理解
// 你可以把它理解为：

// go
// // 这行代码的唯一目的就是：让编译器检查 *clientLogger 是否实现了 io.Writer
// // 如果实现了 → 编译通过，这行代码运行时没有任何效果
// // 如果没实现 → 编译失败，立即发现问题
// 总结
// 这是一种防御性编程手段，在编译期就保证接口实现正确，避免运行时才发现类型不满足接口的问题。在 Go 开源项目中非常常见。

// 这只是在设置log的属性而已，并不是真正的写日志

// Chat
// Thought process
// 用户在问关于 stlog（标准库 log 包）的配置和自定义的 clientLogger.Write 方法之间的关系。让我先读取这个文件来理解完整的上下文。

// \home\guagua\go\go_project\distributed3\log\client.go
// 你的直觉不完全对，实际上这三行代码和 Write 方法是有关系的。让我解释一下：

// 关系链
// log 包的工作流程是这样的：

// plainText
// stlog.Println("hello")
//    ↓
// 1. 拼接前缀: SetPrefix 设置的内容 → "[myservice] - "
// 2. 拼接时间戳: SetFlags 设置的格式 → (Flags=0 则不输出时间)
// 3. 拼接用户消息: "hello\n"
// 4. 调用 Write(data) 把拼接后的完整内容写出
// 所以：

// SetPrefix → 会影响 Write 收到的 data 内容，data 里会包含 "[myservice] - " 这个前缀
// SetFlags(0) → 会影响 Write 收到的 data 内容，设为 0 表示不带时间戳，data 里就不会有时间戳部分
// SetOutput → 决定了最终 Write 方法被谁调用，这里指向了 clientLogger，所以 Write 发出的 HTTP 请求体就是经过前缀和时间戳处理后的完整内容
// 举个具体例子

// go
// // 如果这样调用：
// stlog.Println("hello")

// // SetPrefix("[myservice] - ") + SetFlags(0) 的效果：
// // Write 收到的 data = "[myservice] - hello\n"
// //                        ↑ 前缀           ↑ 用户消息

// // 如果 SetFlags(stlog.LstdFlags)：
// // Write 收到的 data = "[myservice] - 2026/07/02 10:00:00 hello\n"
// //                        ↑ 前缀         ↑ 时间戳              ↑ 用户消息
// 结论
// SetPrefix 和 SetFlags 的设置会决定 Write 收到的 data 的内容格式
// SetOutput 决定了 Write 被谁执行
// 三者共同协作：前两行控制"写什么"，第三行控制"写到哪"
// Write 里的 HTTP 请求发送的就是经过前缀和时间戳格式化后的完整日志内容
