1. 最简 HTTP 服务
go
运行
package main

import (
	"net/http"
)

func main() {
	// 注册路由 + 处理器
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello HTTP"))
	})

	// 启动服务 :8080
	_ = http.ListenAndServe(":8080", nil)
}
2. Handler 两种写法
方式 1：HandleFunc（函数）
go
运行
func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok"))
}

http.HandleFunc("/", handler)
方式 2：实现 http.Handler 接口（对应你前面的 Write 方法）
接口定义：
go
运行
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
示例：
go
运行
type MyHandler struct{}

func (MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("handler impl"))
}

func main() {
	http.Handle("/", MyHandler{})
	_ = http.ListenAndServe(":8080", nil)
}
可以省略接收器变量名，和 Write 写法一致。
3. 获取请求数据
go
运行
func handler(w http.ResponseWriter, r *http.Request) {
	// GET 参数
	name := r.URL.Query().Get("name")

	// POST form 参数
	_ = r.ParseForm()
	age := r.PostForm.Get("age")

	// Path
	path := r.URL.Path
}
4. 返回 JSON
go
运行
import "encoding/json"

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]any{
		"code": 0,
		"msg":  "success",
	}
	_ = json.NewEncoder(w).Encode(data)
}
5. 路由分发（多路复用器 ServeMux）
go
运行
mux := http.NewServeMux()
mux.HandleFunc("/user", userHandler)
mux.HandleFunc("/order", orderHandler)

srv := &http.Server{
	Addr:    ":8080",
	Handler: mux,
}
_ = srv.ListenAndServe()
6. 核心对象说明
http.ResponseWriter w
用来写响应：写内容、状态码、Header。
*http.Request r
包含请求方法、URL、参数、Header、Body。
7. 结合你之前的 io.Writer
日志写入文件可以直接把 fileLog 作为 Writer：
go
运行
type fileLog struct {
	f *os.File
}

func (fl fileLog) Write(p []byte) (n int, err error) {
	return fl.f.Write(p)
}

// 中间件打印请求日志到文件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fl.Write([]byte(r.URL.Path + "\n"))
		next.ServeHTTP(w, r)
	})
}