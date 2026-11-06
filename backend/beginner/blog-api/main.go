package main

import (
	"fmt"
	"net/http"
)

type blogPost struct {
	Id        int      `json:"id"`
	Title     string   `json:"title" validate:"required"`
	Content   string   `json:"content" validate:"required"`
	Category  string   `json:"category" validate:"required"`
	Tags      []string `json:"tags" validate:"required"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

// TODO: Stopped here------
func (bp *blogPost) Create() {

}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
	case "PUT":
	case "DELETE":
	case "GET", "":
	default:
		fmt.Println("Incorrect/unsupported method: ")
		return
	}
}

func main() {
	http.HandleFunc("/posts", postsHandler)
	http.ListenAndServe(":8080", nil)
}
