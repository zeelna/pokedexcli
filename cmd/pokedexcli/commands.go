package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

// Function to handle the 'exit' command.  Usage: Pokedex > exit
func commandExit(conf *ApiConfig, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil

}

// Function to handle the 'help' command.  Usage: Pokedex > help
func commandHelp(conf *ApiConfig, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

// Function to handle the 'map' command to display next 20 LocationAreas.  Usage: Pokedex > map
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

// Function to handle the "mapb' command to go back in Map, i.e previous 20 LocationAreas.  Usage: Pokedex > mapb
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

// Function to handle the "explore <any_location-area>' command to explore Pokemon in the area LocationArea.  Usage: Pokedex > explore pastoria-city-area
func commandExplore(apiConf *ApiConfig, args ...string) error {
	if len(args) <= 0 {
		return fmt.Errorf("not enough arguments to complete command. example usage: explore pastoria-city-area")
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

// Function to handle the "catch <any_pokemon>' command to attempt to catch a Pokemon.  Usage: Pokedex > catch pikachu
func commandCatch(apiConf *ApiConfig, args ...string) error {
	if len(args) <= 0 {
		return fmt.Errorf("not enough arguments to complete command. example usage: catch pikachu")
	}
	pokemonName := args[0]
	url := "https://pokeapi.co/api/v2/pokemon" + "/" + pokemonName

	if _, ok := (*apiConf).Pokedex[pokemonName]; ok {
		return fmt.Errorf("you already have this Pokemon")
	}

	pokemon, err := callAPI[Pokemon](url, apiConf)
	if err != nil {
		return err
	}

	// Calculate if Pokemon catched based on base_experience
	catchChance := rand.IntN(200)
	if pokemon.BaseExperience < catchChance {
		// Successful Catch, add Pokemon to Pokedex
		pokemon.CatchStatus = true
		(*apiConf).Pokedex[pokemon.PokemonName] = pokemon
	} else {
		// Failed Catch
		pokemon.CatchStatus = false
	}

	prefix := fmt.Sprintf("Throwing a Pokeball at %s...", args[0])
	pokemon.print(prefix)
	return nil
}

// Function to handle the "inspect <your_pokemon>' command to inspect the traits of your pokemon.  Usage: Pokedex > inpect pikachu
func commandInspect(apiConf *ApiConfig, args ...string) error {
	if len(args) <= 0 {
		return fmt.Errorf("not enough arguments to complete command. example usage: inspect pikachu")
	}

	pokemonName := args[0]
	pokemon, ok := (*apiConf).Pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}
	prefix := ""
	pokemon.printInspect(prefix)

	return nil
}

// Function to handle the "pokedex' command to display all the Pokemon in your Pokedex you have caught.  Usage: Pokedex > pokedex
func commandPokedex(apiConf *ApiConfig, args ...string) error {
	if 0 >= len((*apiConf).Pokedex) {
		return fmt.Errorf("your have no pokemon in your Pokedex")
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range (*apiConf).Pokedex {
		fmt.Println(" - " + pokemon.PokemonName)
	}
	return nil
}
