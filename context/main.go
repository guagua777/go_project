package main

import (
	"context"
	"fmt"
)

type User struct {
	Name string
}

// Deadline() 方法用于获取 Context 的截止时间，
// Done() 方法用于返回一个只读的 channel，用于通知当前 Context 是否已经被取消，
// Err() 方法用于获取 Context 取消的原因，
// Value() 方法用于获取 Context 中保存的键值对数据。

func main() {
	// 创建闭包， context.Background()
	// 并往闭包中添加数据 context.WithValue()
	ctx := context.WithValue(context.Background(), "name", "guagua")
	// 存储指针
	ctx = context.WithValue(ctx, "user", &User{Name: "guagua"})
	GetUser(ctx)
}

func GetUser(ctx context.Context) {
	// 获取用户名
	fmt.Println(ctx.Value("name"))
	// 类型断言
	// 类型断言的返回值是什么，还是仅仅是做了类型的判断，值还是原来的值
	fmt.Println(ctx.Value("user").(*User).Name)
}
