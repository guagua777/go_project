package http

import (
	"net/http"
	"os"
)

type fileLog struct {
	f *os.File
}

var fl fileLog

func (fl fileLog) Write(p []byte) (n int, err error) {
	return fl.f.Write(p)
}

// 中间件打印请求日志到文件
// 入参为http.Handler, 返回值为http.Handler
// f => g
func loggingMiddleware(next http.Handler) http.Handler {
	// 类型强转
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fl.Write([]byte(r.URL.Path + "\n"))
		next.ServeHTTP(w, r)
	})
}
