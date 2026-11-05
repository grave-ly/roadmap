package main

import (
	// "encoding/json/v2"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Article struct {
	Title   string `json:"title"`
	PubDate string `json:"pubdate"`
	Content []byte `json:"content"`
}

//go:embed templates/*
var fs embed.FS

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	collection, err := populateCollection()
	if err != nil {
		fmt.Println("populating collection in homeHandler: ", err)
	}
	tmpl.Execute(w, collection)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {

}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	collection, err := populateCollection()
	if err != nil {
		fmt.Println("populating collection in adminHandler: ", err)
	}
	tmpl.Execute(w, collection)
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

func populateCollection() ([]Article, error) {
	entries, err := os.ReadDir("articleStore")
	if err != nil {
		return nil, err
	}

	var collection []Article

	for _, entry := range entries {
		path := filepath.Join("articleStore", entry.Name())
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		var newArticle Article
		if err := json.Unmarshal(file, &newArticle); err != nil {
			fmt.Println("unmarshal of article files: ", err)
		}
		collection = append(collection, newArticle)
	}
	return collection, nil
}

func main() {
	fileserver := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			fileserver.ServeHTTP(w, r)
		})))
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/newArtHandler", newArtHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/article", articleHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":8080", nil)
}
