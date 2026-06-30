# go_project


## 参考：
1. https://github.com/xpengkang/go-web-zero
2. 


## go 特性总结
1. 函数：
    1. => 作为普通的值
        1. 可以作为入参和出参，实现高阶函数、统一拦截、日志等
        2. 可以作为某个服务（或者是应用），以参数的形式传递给某个方法
    2. => 接口，本质上就是函数，实现了函数，也就实现了接口
    3. => 方法表达式，即将某个具体类型的方法变为函数，可以实现批量调用、拦截、日志等
    4. => 函数类型，函数类型就是函数，函数类型的调用，就是函数的调用；给函数类型添加方法，添加方法后，如果该方法为某个接口的方法，即实现了该接口，方法中函数类型的调用，就是函数的调用
        1. 链式调用
        2. 实现框架的基石，为什么基于这个就能实现框架？

           什么是框架？

           只有架子，没有实现，那怎么把这个架子搭起来？

           给类型挂载方法，这样就把这个架子搭起来了。
        
           框架的本质：给类型挂载方法，这样这个方法就可以从外面传入，实现框架和具体实现的分离

           有了函数类型后，函数类型挂载方法，用于实现骨架；而具体的方法由外面传入，用于实现具体的业务逻辑
         ```golang
        type HandlerFunc func(ResponseWriter, *Request)

        // ServeHTTP calls f(w, r).
        func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
            f(w, r)
        }
        ```
        1. 以及参考 method-expression-chain-invoke-type-fun 中的代码
2. goroutine
    1. => 本质上为一个continuation，
    2. => continuation + 调度器
3. Context：闭包，解释器的env
    1. => 存值
    2. => 存channel，使用api的形式来完成，用于控制goroutine的执行流程
4. channel
    1. => 输入具体的值
    2. => 输入空结构体 struct{}
    3. => 输出具体的值
    4. => 输出空结构体 struct{}，用作信号
