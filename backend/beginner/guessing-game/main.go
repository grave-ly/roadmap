package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

func guessCheck(guess int, answer int) bool {
	if guess == answer {
		return true
	} else if guess > answer {
		fmt.Printf("Incorrect! The number is less than %d.\n", guess)
		return false
	} else if guess < answer {
		fmt.Printf("Incorrect! The number is greater than %d.\n", guess)
		return false
	}
	return false
}

func textCheck(r *bufio.Reader) int {
	text, _ := r.ReadString('\n')
	text = strings.TrimSpace(text)
	choice, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("string conv to int: ", err)
		return -1
	}
	return choice
}

func main() {
	answer := rand.IntN(100)
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Welcome to the Number Guessing Game!\nI'm thinking of a number between 1 and 100.\nThe amount of chances you have is based on the difficulty you choose.\n\nPlease select the difficulty level:\n1. Easy (10 chances)\n2. Medium (5 chances)\n3. Hard (3 chances)\nEnter your choice: ")

	choice := textCheck(r)
	if choice == -1 {
		return
	}
	if choice < 1 || choice > 3 {
		fmt.Print("Please enter a valid number between 1 and 3 to set the difficulty: ")
		choice = textCheck(r)
	}

	chances := choice
	switch choice {
	case 1:
		fmt.Print("Great! You have selected the Easy difficulty.\nLet's start the game!\n\n")
		chances = 10
	case 2:
		fmt.Print("Great! You have selected the Medium difficulty.\nLet's start the game!\n\n")
		chances = 5
	case 3:
		fmt.Print("Great! You have selected the Hard difficulty.\nLet's start the game!\n\n")
		chances = 3
	}

	for i := 1; i <= chances; i++ {
		fmt.Printf("Enter your guess: ")
		guess := textCheck(r)
		fmt.Print("g: ", guess, "a: ", answer, "\n")
		if guessCheck(guess, answer) {
			fmt.Printf("Congrats! You guessed the correct number in %d attempts\n", i)
			return
		}
	}
}
