package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func addRecord(id int, w *csv.Writer, desc string, amount float64) error {
	now := time.Now()
	date := now.Format("2006-01-02")

	if amount < 0 {
		return errors.New("Amount to be added is not a positive number")
	}
	am := fmt.Sprintf("%.2f", amount)
	idn := strconv.Itoa(id)

	err := w.Write([]string{idn, date, desc, am})
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func summaryRecord(records [][]string, flag *int) {
	var sum float64

	for _, record := range records {
		num, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Println("summary float parse from record:", err)
			continue
		}

		// If monthFlag is nil, add all; otherwise only add matching months
		if flag == nil {
			sum += num
		} else {
			date := record[1]
			month, err := strconv.Atoi(date[5:7])
			if err != nil {
				fmt.Println("summary strconv month from date:", err)
				continue
			}
			if month == *flag {
				sum += num
			}
		}
	}
	if flag == nil {
		fmt.Printf("Total expenses = $%.2f\n", sum)
	} else {
		fmt.Printf("Total expenses in month %d = $%.2f\n", *flag, sum)
	}
}

func deleteRecord(records [][]string, id int) error {
	var temp [][]string
	count := 1
	for _, record := range records {
		rID, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		if rID != id {
			record[0] = strconv.Itoa(count)
			temp = append(temp, record)
			count += 1
		}
	}
	file, err := os.Create("data.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	return w.WriteAll(temp)
}

// args(flags) = add(description, amount), list, summary(month), delete(id)
// store data in a csv file
func main() {
	data, err := os.OpenFile("data.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	r := csv.NewReader(data)
	records, err := r.ReadAll()
	data.Seek(0, io.SeekEnd)
	w := csv.NewWriter(data)

	addFlag := flag.NewFlagSet("add", flag.ExitOnError)
	addDesc := addFlag.String("description", "", "description string")
	addAmount := addFlag.Float64("amount", 0.0, "amount float")

	sumFlag := flag.NewFlagSet("summary", flag.ExitOnError)
	sumMonth := sumFlag.Int("month", -1, "integer month")

	delFlag := flag.NewFlagSet("delete", flag.ExitOnError)
	delID := delFlag.Int("id", -1, "integer id for delete")

	switch os.Args[1] {
	case "add":
		addFlag.Parse(os.Args[2:])
		if *addAmount == 0.0 {
			fmt.Println("Please enter an amount to add")
			return
		}
		err := addRecord(len(records)+1, w, *addDesc, *addAmount)
		if err != nil {
			fmt.Println("error adding record: ", err)
			return
		}
		fmt.Printf("Expense added successfully (ID: %d)\n", len(records)+1)
		return
	case "list":
		fmt.Println("ID\t  Date\t      Description \tAmount")
		for _, record := range records {
			fmt.Printf(" %s\t%s\t%s\t\t$%s\n", record[0], record[1], record[2], record[3])
		}
	case "summary":
		sumFlag.Parse(os.Args[2:])
		if *sumMonth == -1 {
			summaryRecord(records, nil)
		} else {
			summaryRecord(records, sumMonth)
		}
	case "delete":
		delFlag.Parse(os.Args[2:])
		if *delID == -1 {
			fmt.Println("Please enter an ID to delete(-id #)")
			return
		} else {
			err := deleteRecord(records, *delID)
			if err != nil {
				fmt.Println("deletion error: ", err)
			} else {
				fmt.Println("Expense deleted successfully")
			}
		}
		default:
			fmt.Println("Please enter correct commands and flags: -h or --help for more")
			return
	}
}
