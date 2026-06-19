package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 调用适配器处理函数，两个参数，一个http地址，一个是hangler函数
	// r *http.Request 为指针，不需要复制真实的数据，只需要传递指针
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello go"))
	})

	//设置web 服务器,俩个参数，一个监听地址port,一个handler,默认是nil， 采用多路复用mux
	http.ListenAndServe("localhost:8080", nil)
}

func main1() {
	// 普通变量
	var a int = 10
	fmt.Printf("a 的值：%d\n", a)
	fmt.Printf("a 的内存地址：%p\n", &a) // &a 取地址

	// 定义指针变量 *int 代表指向int类型的指针
	var p *int
	p = &a // 把a的地址赋值给指针p
	fmt.Printf("指针p存储的地址：%p\n", p)
	fmt.Printf("*p 解引用，拿到a的值：%d\n", *p)

	// 通过指针修改原变量的值
	*p = 20
	fmt.Printf("修改后 a = %d\n", a) // a 变成20
}
