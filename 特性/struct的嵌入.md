func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}




为什么	return &valueCtx{parent, key, val} 可以当做是Context接口？
21:10
这是一个关于 Go 接口满足（interface satisfaction） 的核心问题。让我结合源码来解释。

关键结构体定义

Go

type valueCtx struct {
    Context       // 匿名嵌入（embedding）父 Context
    key, val any
}
&valueCtx{...} 能当 Context 用，原因有两层：
第一层：valueCtx 嵌入了 Context，所以继承了接口的 4 个方法
Context 接口要求实现 4 个方法：


Go

type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
由于 valueCtx 匿名嵌入了 Context，Go 会自动**委托（delegate）**给内嵌的 Context，也就是说 valueCtx 自动拥有了 Deadline()、Done()、Err() 这 3 个方法——它们直接调用内嵌的 parent 的对应方法。

第二层：valueCtx 自己又**覆盖（override）**了 Value 方法

Go

func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val
    }
    return value(c.Context, key)  // 找不到就委托给父 Context
}
另外还有一个 String() 方法（不影响接口满足，属于额外方法）。

总结：满足接口的完整逻辑
方法	来源
Deadline()	从内嵌的 Context 继承
Done()	从内嵌的 Context 继承
Err()	从内嵌的 Context 继承
Value(key)	valueCtx 自己实现（覆盖了继承的版本）
4 个方法全部具备 → *valueCtx 满足 Context 接口 → &valueCtx{parent, key, val} 可以作为 Context 返回。

本质：Go 的接口是隐式满足的
Go 不像 Java 需要 implements Context，只要一个类型拥有了接口要求的全部方法，它就自动满足该接口。这里 valueCtx 通过"嵌入继承 3 个 + 自己实现 1 个"的方式凑齐了 4 个方法，所以 *valueCtx 就是 Context。