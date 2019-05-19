package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Printf("Listening on localhost:8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
