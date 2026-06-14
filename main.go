package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/zeelna/pokedexcli/internal/pokecache"
)

type ApiConfig struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

type CliCommand struct {
	name        string
	description string
	callback    func(*ApiConfig, ...string) error
}

type LocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

/* // option 2
type LocationAreasResponse struct {
	Count    int                 `json:"count"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []map[string]string		 `json:"results"`
}
*/

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func main() {
	// Cache to store values to optimise API calls
	apiConf := ApiConfig{
		Next:     "https://pokeapi.co/api/v2/location-area", //?offset=0&limit=20",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Second),
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
	}

	// Wait input via command-line, clean the input, get 1st word as the COMMAND and execute .callback(...)
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input += scanner.Text()
		}

		inputSliced := cleanInput(input)
		if len(inputSliced) <= 0 {
			continue
		}
		firstWord := inputSliced[0]

		// verify commands is allowed
		command, ok := commands[firstWord]
		if !ok {
			fmt.Printf("Unknown command: %s\n", firstWord)
			input = "" // do not delete
			continue
		}
		// invoke the command's workflow
		if err := command.callback(&apiConf, inputSliced[1:]...); err != nil {
			fmt.Printf("Could not run: %s\n", firstWord)
			input = "" // do not delete
			continue
		}

		// clean-up, necessary to reset variable. Allows us save a new user-typer command
		input = "" // do not delete
	}

}

func fetchJSON(url string) ([]byte, error) {
	// Create HTTP GET Request. Shorthand //	res, err := http.Get(apiConf.Next)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("could not create request: %w", err)
	}
	// Set Headers
	req.Header.Set("Content-Type", "application/json")

	// Create a Client object, make the HTTP Request and receive the HTTP Response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("network Error: %v", err)
	}
	defer res.Body.Close()

	// Verify successful HTTP GET Request
	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("could not retrieve those location areas. non-OK HTTP status: %s", res.Status)
	}
	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("Response failed with status code: %d and\nbody: %v\n", res.StatusCode, res.Body)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("could not read response body: %w", err)
	}
	return data, nil
}

func callAPI[T any](url string, apiConf *ApiConfig) (T, error) {
	var result T

	// A. If cache has this URL, get cached []byte (cacheEntry exists or not timed-out)
	if cachedData, ok := apiConf.Cache.Get(url); ok {
		err := json.Unmarshal(cachedData, &result)
		return result, err
	}

	// B. If no cache for this URL, call API.
	data, err := fetchJSON(url)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("could not unmarshal: %w", err)
	}

	// Caching successful HTTP Response ([]byte) for 5 seconds each entry
	(*apiConf).Cache.Add(url, data)

	return result, nil
}

// Tightly couple the print function with custom type, to avoid passing wrong parameter / incorrect re-use
/*
type printResponse interface {
	print(string)
}
*/

func (response LocationAreasResponse) print(prefix string) {
	for _, item := range response.Results {
		fmt.Printf("%s\n", item.Name)
	}
}

func (response ExploreResponse) print(prefix string) {
	fmt.Println(prefix)
	for _, item := range response.PokemonEncounters {
		fmt.Printf(" - %s\n", item.Pokemon.Name)
	}
}

func commandExplore(apiConf *ApiConfig, args ...string) error {
	if len(args) <= 0 {
		return fmt.Errorf("not enough arguments to complete command. example: explore pastoria-city-area")
	}
	url := "https://pokeapi.co/api/v2/location-area" + "/" + args[0]

	exploreResponse, err := callAPI[ExploreResponse](url, apiConf)
	if err != nil {
		return err
	}
	prefix := fmt.Sprintf("Exploring %s...\nFound Pokemon:", args[0])
	exploreResponse.print(prefix)
	return nil
}

func commandMap(apiConf *ApiConfig, args ...string) error {
	if (*apiConf).Next == "null" {
		return fmt.Errorf("you're on the last page")
	}
	fullURL := (*apiConf).Next

	pokedexMap, err := callAPI[LocationAreasResponse](fullURL, apiConf)
	if err != nil {
		return err
	}
	(*apiConf).Next = pokedexMap.Next
	(*apiConf).Previous = pokedexMap.Previous

	fmt.Println("Printing next 20 Location Areas!")
	fmt.Printf("HTTP GET Called on: %s\n", (*apiConf).Next)

	prefix := ""
	pokedexMap.print(prefix)
	return nil
}

func commandMapBack(apiConf *ApiConfig, args ...string) error {
	if (*apiConf).Previous == "null" {
		return fmt.Errorf("you're on the first page")
	}
	fullURL := (*apiConf).Previous

	pokedexMap, err := callAPI[LocationAreasResponse](fullURL, apiConf)
	if err != nil {
		return err
	}
	(*apiConf).Next = pokedexMap.Next
	(*apiConf).Previous = pokedexMap.Previous

	fmt.Println("Printing previous 20 Location Areas!")
	fmt.Printf("HTTP GET Called on: %s\n", (*apiConf).Previous)

	prefix := ""
	pokedexMap.print(prefix)
	return nil
}

func commandExit(conf *ApiConfig, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil

}

func commandHelp(conf *ApiConfig, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}
