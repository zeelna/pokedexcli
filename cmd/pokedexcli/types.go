package main

import "github.com/zeelna/pokedexcli/internal/pokecache"

// 1) Store cache to avoid repeated HTTP Request.
// 2) And store Pokedex, "in-memory" map of Pokemon you caught
// 3) Store info from "next: "https://pokeapi.co/api/v2/location-area/?offset=<ANY_X>&limit=20"
// 4) Store "previous":"https://pokeapi.co/api/v2/location-area/?offset=<ANY_X>+20&limit=20?
type ApiConfig struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  map[string]Pokemon
}

// To convert JSON body received from GET https://pokeapi.co/api/v2/location-area/
type LocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// To convert JSON body received from GET https://pokeapi.co/api/v2/location-area/<any_example:pastoria-city-area>
type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// To convert JSON body received from GET https://pokeapi.co/api/v2/pokemon/<any_pokemon__example:pikachu>
type Pokemon struct {
	PokemonName    string `json:"name"`
	Id             int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	CatchStatus bool // populate regardless of API response
}
