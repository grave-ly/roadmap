package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	ID string `json:"id"`
	Type string `json:"type"`
	CreatedAt string `json:"created_at"`
	Actor struct{
		Login string `json:"login"`
	} `json:"actor"`
	Payload struct {
		Action string `json:"action"`
	} `json:"payload"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no github user as arguement")
		return
	}

	arg := os.Args[1]
	template := fmt.Sprintf("https://api.github.com/users/%s/events", arg)

	response, err := http.Get(template)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("body read: ", err)
		return
	}

	var events []Event

	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("unmarshal: ", err)
		return
	}
	for i, event := range events {
		fmt.Printf(" Event %d:\n", i+1)
		fmt.Printf(" Type: %s\n", event.Type)
		fmt.Printf(" Actor: %s\n", event.Actor.Login)
		fmt.Printf(" Action: %s\n", event.Payload.Action)
		fmt.Printf(" Repo: %s\n", event.Repo.Name)
		fmt.Printf(" Created: %s\n\n", event.CreatedAt)
	}
}
