package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pokemastercp/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(time.Second * 5)

type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationArea(url string) (LocationArea, error) {
	var areas LocationArea
	jsonData, cacheHit := cache.Get(url)
	if cacheHit {
		err := json.Unmarshal(jsonData, &areas)
		if err != nil {
			cache.Delete(url)
		} else {
			return areas, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return areas, fmt.Errorf("failed to query PokeAPI: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&areas)
	if err != nil {
		return areas, fmt.Errorf("failed to parse response: %w", err)
	}

	jsonData, err = json.Marshal(areas)
	if err != nil {
		// Failed to marshal JSON data, skipping adding to cache
	} else {
		cache.Add(url, jsonData)
	}
	return areas, nil
}

type LocationEncounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationPokemon(area string) (LocationEncounters, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area
	var locationPokemon LocationEncounters

	jsonData, cacheHit := cache.Get(url)
	if cacheHit {
		err := json.Unmarshal(jsonData, &locationPokemon)
		if err != nil {
			cache.Delete(url)
		} else {
			return locationPokemon, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return locationPokemon, fmt.Errorf("failed to query PokeAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		cache.Add(area, []byte("unknown pokemon"))
		return locationPokemon, fmt.Errorf("unknown pokemon")
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&locationPokemon)
	if err != nil {
		// Failed to marshal JSON data, skipping adding to cache
	} else {
		cache.Add(area, jsonData)
	}

	return locationPokemon, nil
}
