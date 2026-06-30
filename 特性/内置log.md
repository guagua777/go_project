Go 标准库 log（stdlog）简要介绍
1. 包基础
包路径：import "log"
这是 Go 内置的极简日志库，俗称 stdlog，不需要额外安装，开箱即用，适合小工具、临时调试。
核心特点
输出默认写入标准错误 stderr，不是 stdout；
默认每行日志格式：日期 时间 日志内容；
全局只有一个默认 Logger 实例，并发安全；
不支持日志分级（Info/Warn/Error），只有普通打印 + 退出 + 崩溃。
2. 常用全局函数（直接调用）
1）普通打印
go
运行
log.Print("普通日志")
log.Printf("格式化日志：%d", 123)
log.Println("换行日志")
2）打印并退出程序（os.Exit (1)）
go
运行
log.Fatal("发生错误，程序退出")
log.Fatalf("错误码：%v", err)
3）打印并抛出 panic
go
运行
log.Panic("严重错误，触发panic")
3. 自定义 Logger（创建独立日志器）
可以输出到文件，自定义前缀、输出标识：
go
运行
package main

import (
	"log"
	"os"
)

func main() {
	// 参数：输出对象、日志前缀、输出标识
	logger := log.New(os.Stdout, "MYAPP: ", log.LstdFlags)
	logger.Println("自定义日志")
}
日志标识常量（flag）
log.LstdFlags：默认，日期 + 时间
log.Lshortfile：打印文件名 + 行号
log.Lmicroseconds：精确到微秒
组合写法：
go
运行
log.New(w, "prefix: ", log.LstdFlags|log.Lshortfile)
4. 全局 Logger 配置
修改默认 logger 的输出与格式：
go
运行
// 重定向输出到文件
log.SetOutput(fileObj)
// 设置前缀
log.SetPrefix("INFO: ")
// 设置日志flag
log.SetFlags(log.LstdFlags | log.Lshortfile)
5. 优缺点总结
优点
零依赖、代码极简、并发安全；
适合脚本、命令行小工具。
缺点
没有日志级别；
不支持 JSON 格式；
没有切割日志文件、轮转归档能力；
正式项目一般改用：zap、logrus、slog（Go1.21 + 官方新日志库）。
6. 和 slog 的区别
log：老旧极简标准库，无级别；
slog：Go1.21 新增官方结构化日志库，支持级别、JSON、上下文字段，作为 stdlog 的继任者。
需要我给你一份可直接运行的文件日志写入 + 按大小切割的简易封装吗？