package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdat"`
	UpdatedAt   string `json:"updatedat"`
}

var Tasks []Task

func updateTime() string {
	time := time.Now()
	return time.Format("2006-01-02 15:04:05")
}

func nextId() int {
	if len(Tasks) == 0 {
		return 1
	}
	max := 0
	for _, t := range Tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func addTask(description string) {
	nTask := Task{
		ID:          nextId(),
		Description: description,
		Status:      "todo",
		CreatedAt:   updateTime(),
		UpdatedAt:   "",
	}
	Tasks = append(Tasks, nTask)
	fmt.Printf("Task added successfully (ID: %d)\n", nTask.ID)
}

func findIndex(id int) int {
	for i, t := range Tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func updateTask(id int, description string) error {
	i := findIndex(id)
	if i == -1 {
		return errors.New("Task not found when updating")
	}
	Tasks[i].Description = description
	Tasks[i].UpdatedAt = updateTime()
	fmt.Printf("Task (ID: %d) updated \n", id)
	return nil
}

func listTasks(arg string) {
	for _, t := range Tasks {
		updated := t.UpdatedAt
		if updated == "" {
			updated = "no modifications made"
		}
		if len(os.Args) <= 2 {
			fmt.Printf("ID: %d \n Description: %s \n Status: %s \n Created At: %s \n Updated At: %s \n", t.ID, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
		}
		if len(os.Args) > 2 {
			if t.Status == arg {
				fmt.Printf("ID: %d \n Description: %s \n Created At: %s \n Updated At: %s \n", t.ID, t.Description, t.CreatedAt, t.UpdatedAt)
			}
		}
	}
}

func main() {
	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if len(os.Args) < 2 {
		fmt.Println("Please choose an operation")
		return
	}

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var task Task
		if err := decoder.Decode(&task); err != nil {
			fmt.Println("decoding task: ", err)
			return
		}
		Tasks = append(Tasks, task)
	}
	if err = file.Truncate(0); err != nil {
		fmt.Println("truncate after decoding: ", err)
		return
	}
	if _, err = file.Seek(0, 0); err != nil {
		fmt.Println("seek after decoding: ", err)
		return
	}

	var sId int
	sArg := ""
	if os.Args[1] == "list" || os.Args[1] == "add" {
		if len(os.Args) == 3 {
			sArg = os.Args[2]
		}
	} else {
		sId, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("failed to conv arg to int: ", err)
			return
		}
	}

	//commands add, update, delete, mark-in-progress/mark-done, list by status(todo, in-progress, done)
	switch os.Args[1] {
	case "add":
		addTask(sArg)
	case "update":
		if err := updateTask(sId, os.Args[3]); err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		i := findIndex(sId)
		if i == -1 {
			fmt.Println("task not found for delete")
			return
		}
		Tasks = append(Tasks[:i], Tasks[i+1:]...)
		fmt.Printf("Task (ID: %d) deleted\n", sId)

		for idx := range Tasks {
			Tasks[idx].ID = idx + 1
		}
	case "mark-in-progress":
		i := findIndex(sId)
		if i == -1 {
			fmt.Println("task not found for m-i-p")
			return
		}
		Tasks[i].Status = "mark-in-progress"
		Tasks[i].UpdatedAt = updateTime()
		fmt.Printf("marked m-i-p on ID: %d\n", sId)
	case "mark-done":
		i := findIndex(sId)
		if i == -1 {
			fmt.Println("task not found for delete")
			return
		}
		Tasks[i].Status = "done"
		Tasks[i].UpdatedAt = updateTime()
		fmt.Printf("marked done on ID: %d\n", sId)
	case "list":
		listTasks(sArg)

	default:
		fmt.Println("Enter correct arguement(add/update/mark-in-progress/mark-done/list)")
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	for _, t := range Tasks {
		err := encoder.Encode(t)
		if err != nil {
			panic(err)
		}
	}

}
