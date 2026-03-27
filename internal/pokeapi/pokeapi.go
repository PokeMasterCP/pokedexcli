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

	cacheHit, err := cacheCheck(url, &areas)
	if err != nil {
		return areas, err
	} else if cacheHit {
		return areas, nil
	}

	err = queryPokeAPI(url, &areas)
	if err != nil {
		return areas, err
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

	cacheHit, err := cacheCheck(url, &locationPokemon)
	if err != nil {
		return locationPokemon, err
	} else if cacheHit {
		return locationPokemon, nil
	}

	err = queryPokeAPI(url, &locationPokemon)
	if err != nil {
		return locationPokemon, err
	}
	return locationPokemon, nil
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func GetPokemonData(name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	var pokemonData Pokemon

	cacheHit, err := cacheCheck(url, &pokemonData)
	if err != nil {
		return pokemonData, err
	} else if cacheHit {
		return pokemonData, nil
	}

	err = queryPokeAPI(url, &pokemonData)
	if err != nil {
		return pokemonData, err
	}
	return pokemonData, nil

}

func queryPokeAPI[T any](url string, val *T) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to query PokeAPI: %w", err)
	}
	defer resp.Body.Close()

	const notFound = "not found in pokeapi"
	if resp.StatusCode == 404 {
		cache.Add(url, []byte(notFound))
		return fmt.Errorf(notFound)
	} else if resp.StatusCode > 299 {
		return fmt.Errorf("pokeapi responded with %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(val)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	jsonData, err := json.Marshal(*val)
	if err != nil {
		// Failed to marshal JSON data, skipping adding to cache
	} else {
		cache.Add(url, jsonData)
	}

	return nil
}

func cacheCheck[T any](url string, val *T) (bool, error) {
	jsonData, cacheHit := cache.Get(url)
	if cacheHit {
		err := json.Unmarshal(jsonData, val)
		if err != nil {
			cache.Delete(url)
			return false, err
		} else {
			return true, nil
		}
	}
	return false, nil
}
