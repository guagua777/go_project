package main

import (
	"net/http"
	"os"
)

// 统一转换为函数

type Person struct {
	Name string
}

// 给结构体定义方法
func (p Person) Greet(msg string) {
	println(p.Name, msg)
}

// 里面的func(T, string) 对应类型 func(Person, string)
func forEach[T any](slice []T, fn func(T, string)) {
	for _, item := range slice {
		fn(item, "welcome")
	}
}

func main() {
	p := Person{Name: "Tom"}
	p.Greet("hi")

	//------------------------------
	// 方法表达式：接收者变为第一个入参
	methodExpression()

	//------------------------------
	forEachPerson()

	//------------------------------
	// 对 Person.Greet 进行包装
	loggedGreet := wrapLog(Person.Greet)
	loggedGreet(Person{"Jack"}, "test")
}

func forEachPerson() {
	people := []Person{{"A"}, {"B"}, {"C"}}

	// 方法表达式作为通用回调
	forEach(people, Person.Greet)
}

func methodExpression() {
	// 提取为普通函数，接收器变为第一个参数
	greetFunc := Person.Greet
	// 类型：func(Person, string)

	greetFunc(Person{"Alice"}, "hello")
	greetFunc(Person{"Bob"}, "hello")
}

// 本质即为传递方法 wrapLog = f => g
// 在g中实现切面的功能
// 其中f为 (T, string) => (void)，即：f = (T, string) => (void)
// g 为 (T, string) => (void)，即：g = (T, string) => (void)
// func wrapLog(f) => g
// 包装任意接收者为第一个参数的方法
func wrapLog[T any](raw func(T, string)) func(T, string) {
	return func(t T, s string) {
		println("before call")
		raw(t, s)
		println("after call")
	}
}

// 总结：
// 1. 将结构体的方法，转换为普通方法
// 2. 将方法做为参数，使用参数（此处为方法）传递消息

type fileLog struct {
	f *os.File
}

func (fl fileLog) Write(p []byte) (n int, err error) {
	return fl.f.Write(p)
}

var fl fileLog

// 入参为http.Handler，返回值为http.Handler
// 入参为func 出参为func
// 中间件打印请求日志到文件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// 添加切面逻辑
			fl.Write([]byte(r.URL.Path + "\n"))
			next.ServeHTTP(w, r)
		})
}

// 什么是中间件？
// 本质上就是一个函数，在该函数中添加中间件的逻辑
// 中间件 = f => g
// 其中中间件的逻辑，在g中实现
