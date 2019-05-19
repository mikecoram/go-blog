package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Page object
type Page struct {
	Title string
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/post.html")
	t.Execute(w, &Page{Title: "My Title"})
}

func main() {
	http.HandleFunc("/post/", postHandler)
	fmt.Printf("Listening on localhost:8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
