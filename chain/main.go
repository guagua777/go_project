package main

import "fmt"

// func f(g func(string), s string) func(func(string), string) {
// 	return func(h func(string), t string) {
// 		h(g(s))
// 	}
// }

func f(g func(string) string, s string) func(func(string) string, string) string {
	return func(h func(string) string, t string) string {
		// 该函数体是什么？
		v1 := g(s)
		return h(v1 + t)
	}
}

func g(s string) string {
	return s + " world"
}

func h(s string) string {
	return s + " go"
}

func main() {
	// 类似于链式调用
	v := f(g, "hello")(h, " this is param t ")
	fmt.Println(v)
}

// 错误信息是：

// expected '(', found '{'

// 原因
// 问题出在第 9-10 行，函数签名的返回类型换行了：

// Go

// func f(g func(string) string, s string)
// 		func(func(string) string, string) string {
// Go 的词法分析器会自动在行末插入分号。当一行以 ) 结尾时，Go 会自动在后面加上 ;，所以编译器实际看到的是：

// Go

// func f(g func(string) string, s string);    // ← 自动加了分号！
// 		func(func(string) string, string) string {
// 第一行变成了一个不完整的函数声明，编译器期望后面跟 {，却遇到了 func，于是报错。

// 解决方法
// 返回类型必须和函数参数的 ) 在同一行：

// Go

// func f(g func(string) string, s string) func(func(string) string, string) string {
// 	return func(h func(string) string, t string) string {
// 		v1 := g(s)
// 		return h(v1 + t)
// 	}
// }
// Go 自动分号规则
// Go 在以下 token 出现在行末时会自动插入分号：

// 标识符（identifier）
// 各种字面量（整数、浮点数、字符串...）
// break, continue, fallthrough, return
// ++, --
// ), ], }
// 这就是为什么 Go 中 { 必须和函数签名在同一行——否则 ) 末尾会被加分号，导致语法错误。
