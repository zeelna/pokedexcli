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
	callback    func(*ApiConfig) error
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

/*
type LocationAreasResponse struct {
	Count    int                 `json:"count"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []map[string]string		 `json:"results"`
}
*/

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
	}

	// Wait input via command-line, clean the input, get 1st word as the COMMAND and execute .callback(...)
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
			input = "" // do not delete
			continue
		}
		// invoke the command's workflow
		if err := command.callback(&apiConf); err != nil {
			fmt.Printf("Could not run: %s\n", firstWord)
			input = "" // do not delete
			continue
		}

		// clean-up, necessary to reset variable. Allows us save a new user-typer command
		input = "" // do not delete
	}

}

func callAPI(fullURL string, apiConf *ApiConfig) (LocationAreasResponse, error) {
	// Create a HTTP GET Request. Shorthand //	res, err := http.Get(apiConf.Next)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("could not create request: %w", err)
	}
	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	// Create a Client object, make the HTTP Request and receive the HTTP Response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("network Error: %v", err)
	}
	defer res.Body.Close()

	// Verify successful HTTP GET Request
	if res.StatusCode != http.StatusOK {
		return LocationAreasResponse{}, fmt.Errorf("could not retrieve those location areas. non-OK HTTP status: %s", res.Status)
	}
	if res.StatusCode > 299 {
		return LocationAreasResponse{}, fmt.Errorf("Response failed with status code: %d and\nbody: %v\n", res.StatusCode, res.Body)
	}

	// Decode JSON byte-stream into a Go struct
	var resources LocationAreasResponse
	// Option 1: json.Unmarshal works with data that's already in []byte format

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("could not read response body: %w", err)
	}
	if err := json.Unmarshal(data, &resources); err != nil {
		return LocationAreasResponse{}, fmt.Errorf("could not unmarshal the response body: %w", err)
	}
	// Option 2: json.Decoder streams data from io.Reader into a Go struct
	/*
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&resources)
		if err != nil {
			return LocationAreasResponse{}, fmt.Errorf("error decoding response: %v", err)
		}
	*/

	// Location Areas API creates HTTP URLs for
	// 1) next x-amount (20 in our case) Location Areas
	// 2) previous x-amount (20 in our case) Location Areas, there are any (i.e not first
	// Update the next HTTP Request URL, and previous HTTP Request URL
	(*apiConf).Next = resources.Next
	(*apiConf).Previous = resources.Previous

	// Save to Cache after successful HTTP fetch
	(*apiConf).Cache.Add(fullURL, data)
	return resources, nil
}

// Tightly couple the print function with custom type, to avoid passing wrong parameter / incorrect re-use
type printResponse interface {
	printResponse()
}

func (response LocationAreasResponse) printResponse() {
	for _, resultItem := range response.Results {
		fmt.Printf("%s\n", resultItem.Name)
	}
}

// todo: If you're on the first "page" of results, this command should just print "you're on the first page". Example usage:

func commandMap(apiConf *ApiConfig) error {
	fmt.Println("Printing next 20 Location Areas!")
	fmt.Printf("HTTP GET Called on: %s\n", (*apiConf).Next)

	if (*apiConf).Next == "null" {
		return fmt.Errorf("you're on the last page")
	}
	fullURL := (*apiConf).Next

	pokedexMap, err := callAPI(fullURL, apiConf)
	if err != nil {
		return err
	}
	pokedexMap.printResponse()
	return nil
}

func commandMapBack(apiConf *ApiConfig) error {
	fmt.Println("Printing previous 20 Location Areas!")
	fmt.Printf("HTTP GET Called on: %s\n", (*apiConf).Previous)

	if (*apiConf).Previous == "null" {
		return fmt.Errorf("you're on the first page")
	}
	fullURL := (*apiConf).Previous

	pokedexMap, err := callAPI(fullURL, apiConf)
	if err != nil {
		return err
	}
	pokedexMap.printResponse()
	return nil
}

func commandExit(conf *ApiConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil

}

func commandHelp(conf *ApiConfig) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}
