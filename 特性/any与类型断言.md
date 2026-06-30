Go 类型断言（Type Assertion）
一、作用
类型断言用于从空接口 interface{} 中取出原始值，并判断接口变量底层存储的具体类型。
语法有两种写法：
1. 第一种写法（单返回值，出错会 panic）
go
运行
value.(Type)
value：接口变量
Type：要断言的目标类型
示例：
go
运行
package main

import "fmt"

func main() {
	var i interface{} = 10

	// 断言为 int
	n := i.(int)
	fmt.Println(n) // 10

	// 断言失败，直接 panic
	s := i.(string)
	fmt.Println(s)
}
缺点：类型不匹配时程序直接崩溃，一般很少直接用。
2. 第二种写法（双返回值，推荐，无 panic）
go
运行
val, ok := value.(Type)
val：转换后的值
ok：布尔值，true = 类型匹配成功，false = 失败
示例：
go
运行
func main() {
	var i interface{} = "hello"

	s, ok := i.(string)
	if ok {
		fmt.Println("字符串：", s)
	} else {
		fmt.Println("类型不匹配")
	}
}
二、配合 switch：类型分支（type switch）
专门用来批量判断接口的多种类型，语法：
go
运行
switch v := i.(type) {
case int:
	fmt.Println("int类型", v)
case string:
	fmt.Println("string类型", v)
default:
	fmt.Printf("未知类型：%T\n", v)
}
i.(type) 只能用在 switch 里，不能单独使用。
完整例子：
go
运行
package main

import "fmt"

func checkType(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("整数：%d\n", v)
	case string:
		fmt.Printf("字符串：%s\n", v)
	default:
		fmt.Printf("其他类型：%T\n", v)
	}
}

func main() {
	checkType(666)
	checkType("golang")
	checkType(3.14)
}
三、关键知识点
只能对接口变量做断言
普通变量不能直接断言，必须先转为接口。
go
运行
num := 10
// num.(int) // 编译报错，num不是接口
var i interface{} = num
v, ok := i.(int) // 合法
断言只判断动态类型
go
运行
type MyInt int
var i interface{} = MyInt(20)

v, ok := i.(int)    // ok = false
v2, ok2 := i.(MyInt)// ok2 = true
自定义类型和基础类型不相等。
nil 接口断言
接口分为两部分：类型 + 值。只有类型和值都为 nil，接口才等于 nil。
四、总结
基础断言：val, ok := iface.(T)，安全无崩溃。
多类型判断：用 switch iface.(type)。
断言是读取接口底层动态类型的手段，类型转换（conversion）是编译期静态转换，二者不一样。
需要我给你对比一下 类型断言 vs 类型转换 vs 类型判断 的区别吗？




返回值是什么，是仅仅是做了类型的判断，值还是原来的值吗？
一句话结论
第一个返回值 val：是把接口里底层的值，复制一份并转换成目标类型；不是原对象引用，是值副本。
第二个返回值 ok：只做类型判断，匹配为 true，不匹配为 false。
1. 语法拆解
go
运行
val, ok := iface.(T)
ok：纯布尔判断，只检查接口的动态类型是不是 T，不修改任何数据。
val：
类型匹配成功：取出接口包裹的原始值，复制成 T 类型赋值给 val；
类型匹配失败：val 是类型 T 的零值。
示例验证
go
运行
package main
import "fmt"

func main() {
	var i interface{} = 100

	val, ok := i.(int)
	fmt.Println(val, ok) // 100 true

	val2, ok2 := i.(string)
	fmt.Println(val2, ok2) // "" false，val2是string零值
}
2. 值还是原来的值吗？
值类型（int、string、struct）
断言拿到的是副本，修改 val 不会改变接口内部的值：
go
运行
var i interface{} = 10
v, _ := i.(int)
v = 20
fmt.Println(i) // 仍然是10，原数据不变
引用类型（slice、map、指针）
断言得到的依然是副本，但副本存储的是地址，修改指向的内容会影响原数据：
go
运行
type A struct {
	X int
}
var i interface{} = &A{X: 1}

v, _ := i.(*A)
v.X = 999
fmt.Println(i.(*A).X) // 999，内部数据被修改
原因：只是指针变量被复制，指针指向的对象还是同一个。
3. 和单纯类型判断的区别
只用 ok：只判断类型，不提取值。
只要接收了第一个返回值，Go 就会执行值提取 + 类型转换。
如果你只想要判断、不想要取出值，可以用空白标识符：
go
运行
_, ok := i.(int) // 只做判断，不产生变量副本
4. 关键总结
ok：纯粹的类型检测，true/false；
val：成功 = 原始值的同类型副本；失败 = 该类型零值；
不会改写接口变量内部保存的原始值，除非拿到指针再修改指向内容；
类型断言 ≠ 仅判断，它包含：类型检查 + 取值转换两步。
需要我再区分一下：类型断言和 reflect.TypeOf 判断类型的异同吗？


需要
一、核心区别一句话总结
类型断言 i.(T) / i.(type)
编译期语法，专门针对 interface{}；既能判断类型，又能直接把值取出来；性能极高。
反射 reflect.TypeOf(i)
运行期动态类型探测，任意变量都可以检测；只能拿到类型信息，不能直接取出原值；有性能开销。
二、逐条对比
1. 类型断言（Type Assertion）
go
运行
val, ok := i.(int)
适用对象：只能是接口变量
工作时机：编译期就确定目标类型 int，运行时只做一次动态类型比对
能力：
ok：类型是否匹配
val：直接拿到原始数据
失败处理：
单返回值：panic
双返回值：返回零值 + false
优点：零开销，代码简洁
缺点：必须在代码里硬编码要判断的类型，不能动态指定类型
go
运行
// type switch 批量判断
switch v := i.(type) {
case int:
case string:
}
2. 反射 reflect.TypeOf
go
运行
t := reflect.TypeOf(i)
fmt.Println(t.Name())
适用对象：任意变量（普通变量、接口都可以）
工作时机：完全运行期动态解析类型
能力：
只能查询：类型名、种类 Kind、字段、方法等元信息。
想要取出值，必须额外调用 reflect.ValueOf，再用 Interface() 转回接口。
优点：类型可以动态获取，不用提前写死 case，适合通用序列化、ORM、依赖注入
缺点：代码冗长，有运行时性能损耗
示例：
go
运行
var i interface{} = 10
t := reflect.TypeOf(i)
fmt.Println(t.Kind() == reflect.Int) // true

// 想要拿到数值还要多一步
v := reflect.ValueOf(i)
num := v.Int()
三、关键差异对照表
表格
特性	类型断言 .(T)	反射 reflect.TypeOf
使用前提	变量必须是接口类型	任意变量都支持
类型是否硬编码	必须写死类型（int/string），不能动态	运行时动态获取类型，无需硬编码
是否可以直接取值	可以直接提取出原值	只能拿到类型元数据，取值要额外操作
性能	极高，几乎无损耗	有运行时开销
典型场景	已知有限几种类型，分支处理	通用框架、序列化、动态解析
能否自定义类型匹配	严格匹配动态类型，MyInt ≠ int	同样严格区分自定义类型
四、容易踩坑的共同点
二者判断的都是动态类型，不会做隐式转换：
go
运行
type MyInt int
var x interface{} = MyInt(5)

_, ok := x.(int)                  // false
fmt.Println(reflect.TypeOf(x).Kind() == reflect.Int) // true
fmt.Println(reflect.TypeOf(x).Name() == "MyInt")    // true
Kind 拿到底层基础种类
断言与 Type.Name() 匹配的是真实具体类型。
五、使用选择建议
类型只有有限几种（int、string、bool）→ 优先用 type switch 类型断言，简单高效。
写通用库，不知道会传入什么类型 → 只能用 reflect 反射。
能用断言就不要用反射。
需要我再补充：断言、类型转换、reflect 三者放一起的完整对比吗？