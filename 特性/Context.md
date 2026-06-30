1. context就是解释器里面的闭包，就是env，可以提前添加数据
    具体可以看看解释器中的闭包
2. 什么时候添加context参数？
    需要控制goroutine时，就要添加context，用context来控制goroutine的执行
    这跟goroutine之间的通信不同，goroutine之间的通信是通过channel来实现的，而context是用于控制goroutine的执行的


可参考：https://www.fengfengzhidao.com/p/go-context-usage


1. 可以添加值到context中
2. 可以添加cancel操作，用于取消子goroutine的执行