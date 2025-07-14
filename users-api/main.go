package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	fmt.Println("Mock users-api running on port 8080")
	http.ListenAndServe(":8080", nil)
}
