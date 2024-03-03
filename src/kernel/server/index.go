package server

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func StartServer(port int) {
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
