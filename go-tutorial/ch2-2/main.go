package main

import "net/http"

// 自定义handler, 实现ServerHTTP 方法
type myHandler struct {
}

// 自定义handler, 实现ServerHTTP 方法
type helloHandler struct {
}

type aboutHandler struct {
}

// ServeHTTP 不是ServerHTTP
func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home handler"))
}

func (m *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello handler"))
}
func (m *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("about handler"))
}

// 自定义http.handleFunc 函数, 形参和handler函数一样
func welcomeExample(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome!"))
}

func main() {

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: nil,
	}

	mh := myHandler{} // 要使用的是这个指针
	about := aboutHandler{}
	hello := helloHandler{}
	// 不同路径对应不用handler
	http.Handle("/hello", &hello)
	http.Handle("/about", &about)
	http.Handle("/home", &mh)

	http.HandleFunc("/welcome", welcomeExample)
	http.HandleFunc("/welcome2", http.HandlerFunc(welcomeExample))

	server.ListenAndServe()

	// http.ListenAndServe("localhost:8080", nil)
}
