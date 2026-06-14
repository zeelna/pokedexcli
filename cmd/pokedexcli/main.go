package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/zeelna/pokedexcli/internal/pokecache"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*ApiConfig, ...string) error
}

func main() {
	// Cache to store values to optimise API calls
	apiConf := ApiConfig{
		Next:     "https://pokeapi.co/api/v2/location-area", //?offset=0&limit=20",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Second),
		Pokedex:  make(map[string]Pokemon),
	}

	// allowed Command-line commands user can use to receive reply
	commands := map[string]CliCommand{
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
		"map": {
			name:        "map",
			description: "Display Pokedex map of the next 20 Location Areas (API call)",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display Pokedex map of the previous 20 Location Areas (API call)",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore all the Pokemon you can encounter",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon if you can",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View your Pokemon in Pokedex",
			callback:    commandPokedex,
		},
	}

	// Wait input via command-line, clean the input, get 1st word as the COMMAND and execute .callback(...)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		var input string
		fmt.Print("Pokedex > ")
		if scanner.Scan() {

			input += scanner.Text()
		}

		cleaned := cleanInput(input)
		if len(cleaned) <= 0 {
			continue
		}
		commandName := cleaned[0]
		args := cleaned[1:]
		// verify commands is allowed
		command, ok := commands[commandName]
		if !ok {
			fmt.Printf("Unknown command: %s\n", commandName)
			input = "" // do not delete
			continue
		}
		// invoke the command's workflow
		if err := command.callback(&apiConf, args...); err != nil {
			//fmt.Printf("Could not run: %s\nError: %s\n", firstWord, err)
			fmt.Println(err)
			input = "" // do not delete
			continue
		}

		// clean-up, necessary to reset variable. Allows us save a new user-typer command
		input = "" // do not delete
	}

}
