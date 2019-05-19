package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Page object
type Page struct {
	Title string
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/post.html")
	t.Execute(w, &Page{Title: "My Title"})
}

func getDbConnection() (db *sql.DB) {
	db, err := sql.Open("postgres", "user=postgres password=secret dbname=go-blog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	getDbConnection()
	http.HandleFunc("/post/", postHandler)
	fmt.Printf("Listening on localhost:8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
