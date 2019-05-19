package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// Post object
type Post struct {
	Title   string
	Content string
	URLSlug string
}

// HomeData object
type HomeData struct {
	Posts []Post
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	db := getDbConnection()
	rows, err := db.Query("SELECT title, url_slug FROM posts")
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	}
	var posts []Post
	for rows.Next() {
		var (
			title   string
			urlSlug string
		)
		rows.Scan(&title, &urlSlug)
		posts = append(posts, Post{Title: title, URLSlug: urlSlug})
	}
	t, _ := template.ParseFiles("views/index.html")
	t.Execute(w, HomeData{Posts: posts})
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		body, _ := ioutil.ReadFile("views/new_post.html")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write(body)
	} else if r.Method == "POST" {
		r.ParseForm()
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		slug := strings.ReplaceAll(strings.ToLower(title), " ", "-")

		db := getDbConnection()
		_, err := db.Query(`
			INSERT INTO posts
			(title, content, url_slug, created_on)
			VALUES
			($1, $2, $3, NOW());
		`, title, content, slug)
		if err == nil {
			http.Redirect(w, r, "/home", http.StatusFound)
		} else {
			http.Error(w, "Failed to save the post.", http.StatusInternalServerError)
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.Split(strings.ToLower(r.URL.Path), "/post/")[1]

	db := getDbConnection()
	rows, err := db.Query("SELECT title, content FROM posts WHERE url_slug = $1", slug)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	}

	if rows.Next() {
		var (
			title   string
			content string
		)
		rows.Scan(&title, &content)
		t, _ := template.ParseFiles("views/post.html")
		t.Execute(w, &Post{Title: title, Content: content})
	} else {
		http.NotFound(w, r)
	}
}

func getDbConnection() (db *sql.DB) {
	db, err := sql.Open("postgres", `
		host=127.0.0.1
		port=5432
		user=postgres
		password=secret
		dbname=blog
		sslmode=disable
	`)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/new-post/", newPostHandler)
	http.HandleFunc("/post/", postHandler)

	fmt.Printf("Listening on localhost:8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
