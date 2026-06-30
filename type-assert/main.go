package main

import "fmt"

type MyInt int

func checkType(i interface{}) {
	// 这里面的type，就是type关键字
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
