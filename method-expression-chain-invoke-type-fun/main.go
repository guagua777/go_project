package main

// 给方法类型添加方法

// 1. 定义方法类型
type g func(string)

// 2. 给方法类型添加方法
// 可以添加任意类型的方法，也可以有返回值
func (g g) Greet(msg string) {
	// println(msg)
	g(msg)
}

// 添加的方法，参数不需要和方法类型一致
func (g g) Greet2(msg string, a int) {
	println(msg, a)
	g(msg)
}

func (g g) Greet_Chain(msg string) g {
	// println(msg)
	g(msg)
	return g
}

func (g g) Greet_Chain2(msg string) g {
	return func(str string) {
		println(msg)
		g(str)
	}
}

func main() {
	gFun := g(func(str string) {
		println("first", str)
	})
	gFun.Greet("hello")
	gFun.Greet_Chain("hello").Greet("world")
	gFun.Greet_Chain("hello")("world")
	gFun.Greet_Chain2("hello").Greet("world")
}

func main1() {
	gFun := func(str string) {
		println(str)
	}
	gGreet := g.Greet
	gGreet(gFun, "hello")

	gFunNamed := g(gFun)

	gFunNamed.Greet("hello")

	// 问题在第 31 行：

	// go
	// r := gFunNamed.Greet("hello")
	// Greet 方法没有返回值：

	// go
	// func (g g) Greet(msg string) {   // 没有返回类型
	// 但你试图把它的返回值赋给 r，所以编译器报错：(no value) used as value——没有返回值的东西被当成值来用了。

	// 解决方法
	// 取决于你的意图：

	// 如果只是想调用 Greet，不需要返回值，删掉赋值即可：

	// go
	// gFunNamed.Greet("hello")
	// 如果想支持链式调用，让 Greet 返回 g 类型本身：

	// go
	// func (g g) Greet(msg string) g {
	// 	g(msg)
	// 	return g
	// }
	// 这样就能链式调用了：

	// go
	// gFunNamed.Greet("hello").Greet("world")

	// r := gFunNamed.Greet("hello")
	// println(r)
}
