package main

import (
	"encoding/json"
	"net/http"

	"github.com/guagua777/web-tutorial/ch23/middleware"
)

type Company struct {
	ID      int
	Name    string
	Country string
}

func main() {
	http.HandleFunc("/companies", func(w http.ResponseWriter, r *http.Request) {
		c := Company{
			ID:      123,
			Name:    "gggoolle",
			Country: "USA",
		}
		enc := json.NewEncoder(w)
		enc.Encode(c)
	})
	// 使用中间件
	http.ListenAndServe("localhost:8080", new(middleware.AuthMiddleware))
}
