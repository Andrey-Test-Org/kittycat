package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting orders service")

	http.HandleFunc("/orders", HandleOrders)
	http.HandleFunc("/order", HandleOrder)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	http.ListenAndServe(":9090", nil)
}
