package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	/*
		fmt.Println("-----------------")
		fmt.Printf("%v\n", cleanInput("Hello, World!"))
		fmt.Printf("%v\n", cleanInput("pikacHu AsH "))
		fmt.Printf("%v\n", cleanInput("   "))
		fmt.Printf("%v\n", cleanInput(""))
		fmt.Println("-----------------")
	*/
	type cliCommand struct {
		name        string
		description string
		callback    func() error
	}

	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Instructions for Pokedex",
			callback:    commandHelp,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input += scanner.Text()
		}

		textSlice := cleanInput(input)
		if len(textSlice) <= 0 {
			continue
		}
		firstWord := textSlice[0]

		// verify commands is allowed
		command, ok := commands[firstWord]
		if !ok {
			fmt.Printf("Unknown command: %s\n", firstWord)
			continue
		}

		if err := command.callback(); err != nil {
			fmt.Printf("Unknown command: %s\n", firstWord)
			continue
		}
		// invoke the command's workflow

		input = "" // clean-up, necessary to reset variable. Allows us save a new user-typer command
	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil

}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}
