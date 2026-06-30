这里面的fileLog(destination)是强转还是构造​

这是类型转换（强转），不是构造函数调用。

在 Go 中，fileLog(destination) 的含义取决于 fileLog 的定义。看代码第 11 行：


go
type fileLog string
fileLog 是基于 string 的自定义类型（named type），而 destination 的类型是 string。所以 fileLog(destination) 就是把 string 类型的值转换为 fileLog 类型，语法形式是 T(v) —— 这是 Go 的标准类型转换语法。

两者的区别
类型转换 T(v)	构造函数
Go 语法	fileLog(destination)	Go 没有构造函数语法
本质	将底层类型兼容的值转为目标类型	惯例上用 NewXxx() 函数模拟
例子	int(3.14), fileLog("log.txt")	NewFileLog("log.txt")
为什么这里需要转换？
因为 fileLog 实现了 io.Writer 接口（Write 方法），而 stdlog.New 的第一个参数需要 io.Writer。直接传 destination（string 类型）不行，必须先转为 fileLog 类型，这样 Go 才能识别它满足 io.Writer 接口。

在 Write 方法内部，又通过 string(fl) 把 fileLog 转回了 string，用于 os.OpenFile 打开文件。



type 新类型名 底层类型

// 基于int创建一个全新类型MyInt
type MyInt int

func main() {
    var a MyInt = 10
    var b int = 20
    // a = b 编译报错！MyInt 和 int 是两种不同类型，不能直接赋值
    a = MyInt(b) // 必须显式转换
}

底层存储和原类型完全一致；
属于完全不同的类型，不支持隐式类型转换；
可以给新类型绑定方法。

type MyInt int

// 为MyInt增加方法
func (m MyInt) Double() int {
    return int(m) * 2
}
原生 int 没有这个方法，只有 MyInt 拥有。


三、类型别名：type 别名 = 原类型
加 = 是别名，二者完全等价，视为同一个类型。
type MyInt = int

func main() {
    var a MyInt = 10
    var b int = a // 直接赋值，无报错
}

区别总结：
表格
写法	类型关系
type T int	新类型，不能隐式转换
type T = int	类型别名，完全等同，可混用
四、定义结构体 type struct
go
运行
type User struct {
    ID   int
    Name string
}
type + struct 定义结构体类型。
五、定义接口 type interface
go
运行
type Reader interface {
    Read() ([]byte, error)
}
一句话总结
type T X：创建新类型，独立类型，可加方法；
type T = X：创建类型别名，仅仅是换个名字，类型不变；
type T struct{}：定义结构体；
type T interface{}：定义接口。



1. 只能给自定义类型加方法
只有用 type 新建出来的类型，才能绑定方法。
原生类型（int、string）不能直接加方法。
写法格式
go
运行
func (接收者 变量 类型) 方法名(参数) 返回值 {
}
值接收者（拷贝，不修改原变量）
go
运行
package main

import "fmt"

// 新建自定义类型
type Num int

// 值接收者
func (n Num) AddTwo() Num {
    return n + 2
}

func main() {
    var x Num = 5
    fmt.Println(x.AddTwo()) // 7
}
指针接收者（修改原值，共享内存）
go
运行
type Num int

// 指针接收者
func (n *Num) Inc() {
    *n++
}

func main() {
    var x Num = 1
    x.Inc()
    fmt.Println(x) // 2
}
2. 结构体上定义方法（最常用）
go
运行
type User struct {
    Name string
    Age  int
}

// 值接收者
func (u User) Hello() string {
    return "hello " + u.Name
}

// 指针接收者，修改成员
func (u *User) Grow() {
    u.Age++
}
3. 关键规则
值接收者：操作副本，外部变量不会改变；
指针接收者：操作原对象，可以修改字段；
Go 自动兼容调用：
值变量可以调用指针方法
指针变量可以调用值方法
不能给别名类型额外加方法：
go
运行
type A = int   // 别名，和int是同一个类型
// 无法给A定义新方法，编译报错
4. 值接收者 vs 指针接收者选择口诀
只读取数据 → 值接收者
需要修改字段 → 指针接收者
结构体很大，减少拷贝开销 → 优先指针接收者



函数类型定义：type + func
语法：
go
运行
type 类型名 func(入参) 返回值
1. 基础示例
go
运行
// 定义一个函数类型
type Handler func(string) error
Handler 不再是普通变量，是一种自定义函数类型；
凡是签名为 func(string) error 的函数，都属于该类型。
使用：
go
运行
func myHandle(s string) error {
    return nil
}

func main() {
    var h Handler = myHandle
    h("test")
}
2. 给函数类型绑定方法（重点）
普通函数不能加方法，但 type 出来的函数类型可以。
go
运行
type Handler func(string) error

// 给函数类型添加方法
func (h Handler) WrapLog(s string) error {
    println("start")
    return h(s)
}
调用：
go
运行
var h Handler = myHandle
h.WrapLog("hello")
这是中间件、回调包装最经典的写法。
3. 带多参数、多返回值
go
运行
type Calc func(a, b int) int
4. 和普通函数变量的区别
var f func(int) int：只是一个函数变量；
type F func(int) int：创建了一个新类型，可以附加方法。
5. 实战：HTTP 处理器原型
go
运行
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)
}
这就是标准库 http.HandlerFunc 的源码原理。





-------
// 定义一个函数类型
type Handler func(string) error
Handler 不再是普通变量，是一种自定义函数类型；
凡是签名为 func(string) error 的函数，都属于该类型。

同理：

// 把函数签名升级为【具名自定义类型】
type HandlerFunc func(w ResponseWriter, r *Request)

// 给这个函数类型补上接口要求的方法
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r) // 内部调用原函数
}

凡是签名为 func(w ResponseWriter, r *Request) 的函数，都属于HandlerFunc类型 
->
而HandlerFunc类型有ServeHTTP方法，所以凡是签名为 func(w ResponseWriter, r *Request) 的函数都所有 ServeHTTP方法
-> ServeHTTP方法 又是 http.Handler 接口的方法
所以凡是签名为 func(w ResponseWriter, r *Request) 的函数都实现了 http.Handler 接口。


                               外部方法，签名相同
                                    |
                                    |
                                    |
HandlerFunc 类型  --------->    func(...) 方法签名
    |                                                                                                |                                                                       接口 
    |                                                                        |
    |                                                                        |
    |-----------------------------------------》 类型的方法a   ----------->  方法a


