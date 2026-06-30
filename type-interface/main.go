package main

type Inter interface {
	Call(s string)
}

func (i Inter) Call2(s string) {
	i.Call(s)
}

func main() {

}

// 错误信息是：

// invalid receiver type Inter (pointer or interface type)

// 原因
// 第 8 行 func (i Inter) Call2(s string) 中，你试图用接口类型 Inter 作为方法的接收者（receiver），但 Go 不允许接口类型作为接收者。

// 方法接收者必须是具体类型（或具体类型的指针），不能是接口或指针类型的别名。

// 为什么不允许？
// 接口是行为的抽象描述，它定义"能做什么"，而不是"是什么"。方法必须绑定在具体的实体上，而不是抽象的描述上。如果允许接口做接收者，就会产生逻辑矛盾——接口本身没有数据和状态，你无法对一个抽象概念调用方法。

// 正确做法
// 如果你想让所有实现了 Inter 接口的类型都拥有 Call2 方法，有两种常见方式：

// 方式1：把 Call2 写成普通函数，参数为接口类型

// Go

// func Call2(i Inter, s string) {
//     i.Call(s)
// }
// 方式2：定义一个具体类型来嵌入接口（较少见）

// Go

// type Wrapper struct {
//     Inter
// }

// func (w Wrapper) Call2(s string) {
//     w.Call(s)
// }
// 总结
// 写法	是否合法
// func (i Inter) Call2()	❌ 接口不能做接收者
// func (i *Inter) Call2()	❌ 指向接口的指针也不能做接收者
// func (i MyStruct) Call2()	✅ 具体类型可以做接收者
// func Call2(i Inter)	✅ 接口可以作为普通函数参数
