package main

import "fmt"

func (response Pokemon) print(prefix string) {
	fmt.Println(prefix)
	if response.CatchStatus {
		msg := fmt.Sprintf("%s was caught!", response.PokemonName)
		fmt.Println(msg)
	} else {
		msg := fmt.Sprintf("%s escaped!", response.PokemonName)
		fmt.Println(msg)
	}
}

func (response Pokemon) printInspect(prefix string) {
	fmt.Printf(prefix)
	fmt.Printf("Name: %s\n", response.PokemonName)
	fmt.Printf("Height: %d\n", response.Height)
	fmt.Printf("Weight: %d\n", response.Weight)

	fmt.Println("Stats:")
	for _, item := range response.Stats {
		fmt.Printf("  - %s:  %d\n", item.Stat.Name, item.BaseStat)
	}

	fmt.Println("Types:")
	for _, item := range response.Types {
		fmt.Println("  - " + item.Type.Name)
	}
}

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
