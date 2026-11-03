package main

import (
	// "encoding/json/v2"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Article struct {
	Title   string `json:"title"`
	PubDate string `json:"pubdate"`
	Content []byte `json:"content"`
}

var Collection []Article

//go:embed templates/*
var fs embed.FS

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, Collection)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {

}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, Collection)
}

func editHandler(w http.ResponseWriter, r *http.Request) {

}

func newHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data, _ := fs.ReadFile("templates/new.html")
	w.Write(data)
}

func newArtHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	art := Article{
		Title:   r.FormValue("articleTitle"),
		PubDate: r.FormValue("pubDate"),
		Content: []byte(r.FormValue("content")),
	}
	data, err := json.Marshal(art)
	if err != nil {
		fmt.Println("new article ", err)
		return
	}
	entries, err := os.ReadDir("./articleStore/")
	if err != nil {
		fmt.Println("reading dir in new art ", err)
	}
	path := fmt.Sprintf("articleStore/%d.json", len(entries)+1)
	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Println("writing new article to files", err)
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// func populateCollection() error {
// 	entries, err := os.ReadDir("./articleStore/")
// 	if err != nil {
// 		return err
// 	}
// 	for _, entry := range entries {
// 		file, err := os.ReadFile(entry.Name())
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		var newArticle Article
// 		json.Unmarshal(file, &newArticle)
// 	}
// 	return nil
// }

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/newArtHandler", newArtHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/article", articleHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":8080", nil)
}
