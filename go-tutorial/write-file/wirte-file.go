package main

import (
	"os"
)

func main2() {
	// 创建切片
	content := []byte("hello go file\n")
	// 覆盖写入，不存在则创建
	err := os.WriteFile("test.txt", content, 0644)
	if err != nil {
		// panic
		panic(err)
	}
}

type T struct {
	Name string
}

func (t T) Hello(s string) {
	println(s + t.Name)
}

func main() {
	// 方法表达式：接收者变为第一个入参
	f := T.Hello
	// f 的类型：func(T, string)

	var t T
	f(t, "hi") // 等价于 t.Hello("hi")
}
