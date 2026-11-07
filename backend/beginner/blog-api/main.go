package main

import (
	"context"
	"database/sql"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"

	_ "modernc.org/sqlite"
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

// TODO: Stopped here------add blog to db and add timestamps
func (bp *blogPost) Create(db *sql.DB) error {

	return nil
}

func postsHandler(db *sql.DB) http.HandlerFunc {
	var entry blogPost
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("reading data in postsHandler: ", err)
		}
		defer r.Body.Close()

		switch r.Method {
		case "POST":
			if err := json.Unmarshal(data, &entry); err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				fmt.Println("post Unmarshal: ", err)
				return
			}
			if err = entry.Create(db); err != nil {
				http.Error(w, "Failed to create blog post", http.StatusInternalServerError)
				fmt.Println("create blog in db: ", err)
				return
			}
		case "PUT":
		case "DELETE":
		case "GET", "":
		default:
			http.Error(w, "Incorrect/unsupported method", http.StatusMethodNotAllowed)
			fmt.Println("Incorrect/unsupported method: ")
			return
		}
	}
}

func main() {

	ctx := context.Background()

	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		fmt.Println("db failed to open: ", err)
	}

	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		fmt.Println("failed to ping db: ", err)
		return
	}

	http.HandleFunc("/posts", postsHandler(db))
	http.ListenAndServe(":8080", nil)
}
