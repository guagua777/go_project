package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "filesystem"+r.URL.Path)
	})

	http.ListenAndServe("localhost:8080", nil)
}
