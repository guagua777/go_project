package sub

import "net/http"

type MyHandler struct{}

// 指针方法
func (*MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

// 报错：MyHandler{} 没有实现 http.Handler
var _ http.Handler = MyHandler{}

// 正常：&MyHandler{} 实现了接口
// &取地址
var _ http.Handler = &MyHandler{}
