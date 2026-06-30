package valuepointer

import "net/http"

type MyHandler struct{}

func (MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("handler impl"))
}

func main() {
	http.Handle("/", MyHandler{})
	_ = http.ListenAndServe(":8080", nil)
}
