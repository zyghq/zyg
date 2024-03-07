package main

import (
	"net/http"
)

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// func main() {
// 	fmt.Println("Hello, World!")
// }
