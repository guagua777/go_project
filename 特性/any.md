一、结论先行
在 Go 1.18+：
go
运行
type any = interface{}
any 只是 interface{} 的别名，二者完全等价，没有任何区别。
二、interface {} 空接口
1. 定义
go
运行
interface{}
一个没有定义任何方法的接口。
2. 核心特性
Go 中所有类型都自动实现了空接口：
int、string、bool
结构体、切片、map、函数、chan
自定义类型
所以：
go
运行
var x interface{}
x = 123
x = "hello"
x = []int{1,2,3}
空接口变量可以存放任意类型的值。
3. 底层结构（runtime.eface）
空接口变量包含两部分：
类型指针（_type）：记录存入值的真实类型
数据指针：指向值的内存
非空接口用 iface，空接口用 eface。
4. 取值必须类型断言
不能直接把 interface {} 赋值给具体类型：
go
运行
var v interface{} = 10
a := v.(int) // 类型断言
三、any 是什么
源码（builtin.go）
go
运行
// any is an alias for empty interface
type any = interface{}
any = interface{}
编译后完全一样，字节码无差别
只是语法糖，让代码更简洁易读
对比写法：
go
运行
func foo(a interface{}) {}
func bar(a any) {} // 和上面一模一样
什么时候用 any？
泛型代码（Go1.18 + 泛型里首选 any）
go
运行
func print[T any](val T) {}
通用容器、万能参数，可读性比光秃秃的 interface{} 更好。
四、关键区别总结
表格
写法	本质	适用版本
interface{}	原生空接口类型	所有 Go 版本
any	interface {} 的类型别名	Go 1.18+
行为、内存布局、类型断言、装箱拆箱逻辑完全一致。
五、常见误区
❌ 误区：any 是新类型，性能不一样
✅ 正解：只是别名，编译器直接替换成 interface {}，无性能差异。
❌ 误区：any 只能用于泛型
✅ 正解：普通变量、函数参数都能随便用。
go
运行
var data any = "test"
❌ 把 any 当成 object
Go 没有 object，any/interface {} 只是静态语言里的 “任意类型容器”，依然是静态类型，运行时要做类型检查。
六、补充：空接口带来的问题
失去编译期类型检查，容易运行 panic
会发生值装箱，产生堆内存分配（escape to heap）
频繁使用 interface {} 会影响性能，优先使用泛型代替万能空接口。





一、先讲：eface（空接口）和 iface（非空接口）
1）空接口：interface{} → 底层结构体叫 eface
空接口没有任何方法，只需要保存两样东西：
实际值的类型信息（_type）
实际数据的地址
go
运行
// runtime 源码简化版
type eface struct {
    _type *type  // 存类型：int / string / struct...
    data  unsafe.Pointer // 存数据指针
}
当你写：
go
运行
var x interface{} = 10
底层就构造了一个 eface 对象，记录：类型是 int，数据是 10。
2）非空接口（带方法）→ iface
举个例子：
go
运行
type Reader interface {
    Read()
}
这个接口有方法，底层结构是 iface：
go
运行
type iface struct {
    tab  *itab   // 类型 + 方法表（记录该类型实现了接口里哪些方法）
    data unsafe.Pointer
}
区别：
空接口 eface：只存类型 + 数据，没有方法表
非空接口 iface：多了一张方法表 itab，用来校验类型是否实现接口
一句话总结：
只要接口里有方法 → iface；
接口里空空如也 interface{} → eface。
二、第二块：为什么必须做类型断言 v.(int)？
1. 根源：静态类型限制
go
运行
var v interface{} = 10
// var a int = v   // 这一行直接编译报错
原因：
变量 v 的静态类型永远是 interface {}，编译器只知道它是空接口，并不知道里面装的是 int。
Go 是静态强类型语言，不允许直接把接口变量自动转为具体类型，必须手动告诉编译器：“我确定里面是 int”。
2. 类型断言语法
写法 1（断言失败直接 panic）
go
运行
a := v.(int)
写法 2（安全写法，不会崩溃）
go
运行
a, ok := v.(int)
if ok {
    // 转换成功
}
3. 运行时发生了什么？
程序取出 eface 里面保存的 _type，对比是不是 int：
匹配成功：把 data 里的值拿出来赋值给 a
不匹配：直接崩溃（不加 ok 的写法）
三、举个直观例子帮你吃透
go
运行
package main

func main() {
    var x interface{}
    x = 666    // 打包成 eface{type:int, data:&666}

    // 不能直接 b := x
    b := x.(int) // 拆开eface，校验类型，取出数值
    println(b)
}
四、容易混淆的点补充
any 本质还是 interface{}，底层依然是 eface，不会变成 iface。
只有带方法的自定义接口才会使用 iface 结构。
类型断言是运行时检查，不是编译检查，所以滥用空接口容易出现线上 panic。