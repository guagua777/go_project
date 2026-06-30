package main

func f() int {
	i := 0
	defer func() { i++; println("i is", i) }()
	return i // 先把返回值设为0，再执行defer i变成1，最终返回0
}

func main() {
	v := f()
	println(v)
}
