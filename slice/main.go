package main

import "fmt"

func modify(s []int) []int {
	s[0] = 999 // 外部会变
	// append如果不扩容，不会改变切片的底层数组
	s = append(s, 100) // 扩容后内部是新切片，外部无变化
	return s
}

func main() {
	s := []int{1, 2, 3}
	// 打印切片内容
	fmt.Printf("%+v\n", s) // 结构体时有用，切片和 %v 一样
	// r := modify(s)
	r := modify(s[:1])
	fmt.Printf("%+v\n", r) // 结构体时有用，切片和 %v 一样
	fmt.Printf("%+v\n", s)

}
