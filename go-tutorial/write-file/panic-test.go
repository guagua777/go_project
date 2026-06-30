package main

func test() {
	defer func() {
		// 捕获崩溃
		if err := recover(); err != nil {
			println("捕获异常：", err)
		}
	}()

	panic("严重错误")
	println("这里永远不会执行")
}

func main1() {
	test()
	println("主程序继续运行")
}
