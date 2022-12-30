package main

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	if err := http.ListenAndServe(":9100", nil); err != nil {
		panic(err)
	}
}
