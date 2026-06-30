package main

import "fmt"

// 1. 定义具名函数类型（关键：只有具名类型才能绑定方法）
type Handler func()

// 2. 链式中间件1：日志
func (h Handler) WithLog() Handler {
	return func() {
		fmt.Println("[日志] 开始执行")
		h()
		fmt.Println("[日志] 执行结束")
	}
}

// 3. 链式中间件2：鉴权
func (h Handler) WithAuth() Handler {
	return func() {
		fmt.Println("[鉴权] 校验权限")
		h()
	}
}

// 4. 链式中间件3：超时处理
func (h Handler) WithTimeout() Handler {
	return func() {
		fmt.Println("[超时] 开启超时控制")
		h()
	}
}

// 业务主逻辑
func business() {
	fmt.Println("执行业务代码")
}

func main() {
	// 裸函数强制转为具名类型
	chain := Handler(business).
		WithLog().
		WithAuth().
		WithTimeout()

	// 执行整条链路
	chain()
}
