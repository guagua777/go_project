package log

import (
	"io"
	stdlog "log"
	"net/http"
	"os"
)

/**
整体内容：
1. 写log
2. 注册服务，即注册handler
**/

var log *stdlog.Logger

type fileLog string

// 写日志到文件中
func (fl fileLog) Write(data []byte) (int, error) {
	// 创建
	// 只可写
	// 追加
	// 权限 linux中的文件权限
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

// 运行自定义的log的时候，想将log指向某一个具体的文件地址
// 服务端log的前缀设置
// 给log赋值
func Run(destination string) {
	// 因为fileLog实现了io.Writer接口，所以可以将fileLog作为参数传递给stdlog.New函数
	log = stdlog.New(fileLog(destination), "go: ", stdlog.LstdFlags) // LstdFlags表示使用标准的日期和时间格式（2006/01/02 15:04:05）来记录日志
}

// 为什么这个可以作为一个服务？
// 因为是一个方法，可执行
// 可执行，就可以作为一个服务
// 注册log的http handler
func RegisterHandler() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost: // post请求
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
			// 不用break，go中自带
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v", message)
}
