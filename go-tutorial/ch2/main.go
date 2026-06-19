package main

import "net/http"

// 自定义handler, 实现ServerHTTP 方法
type myHandler struct {
}

// ServeHTTP 不是ServerHTTP
func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home handler"))
}

func main() {

	// server := http.Server{
	// 	Addr:    "localhost:8080",
	// 	Handler: nil,
	// }
	// server.ListenAndServe()

	mh := myHandler{} // 要使用的是这个指针
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: &mh,
	}
	server.ListenAndServe()

	// http.ListenAndServe("localhost:8080", nil)
}
