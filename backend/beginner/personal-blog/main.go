package main

import (
	"encoding/json/v2"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Article struct {
	Title   string `json:"title"`
	PubDate string `json:"pubdate"`
	Content string `json:"content"`
}

var Collection []Article

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))

	tmpl.Execute(w, Collection)

}

func articleHandler(w http.ResponseWriter, r *http.Request) {

}

func adminHandler(w http.ResponseWriter, r *http.Request) {

}

func editHandler(w http.ResponseWriter, r *http.Request) {

}

func newHandler(w http.ResponseWriter, r *http.Request) {

}

func populateCollection() error {
	entries, err := os.ReadDir("./articleStore/")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		file, err := os.ReadFile(entry.Name())
		if err != nil {
			panic(err)
		}
		var newArticle Article
		json.Unmarshal(file, &newArticle)
	}
	return nil
}

func main() {
	if err := populateCollection(); err != nil {
		fmt.Println(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/article", articleHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":8080", nil)
}
