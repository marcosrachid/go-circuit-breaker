package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", hello)
	log.Printf("Serving at localhost:8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
