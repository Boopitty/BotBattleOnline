package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Boopitty/BotBattleOnline/internal/gamelogic"
)

func main() {
	fmt.Println("State your name:")
	scanner := bufio.NewScanner(os.Stdin)
	defer os.Stdin.Close()

	var input string
	if scanner.Scan() {
		input = scanner.Text()
	}
	user := gamelogic.NewPlayer(input)

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v", err)
	}

	for {
		fmt.Print("> ")
		if scanner.Scan() {
			input := scanner.Text()
			switch input {
			case "profile":
				fmt.Printf("Here is your info:\n%v\n", *user)
			case "quit":
				fmt.Println("Have a nice day!")
				os.Exit(0)
			default:
				fmt.Printf("Your order was: %s\n", input)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error: %v", err)
		}
	}
}
