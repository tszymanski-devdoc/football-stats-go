package main

import (
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from HTTP GET"))
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Football stats data"))
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/stats", statsHandler)
	http.ListenAndServe(":8080", nil)
}
