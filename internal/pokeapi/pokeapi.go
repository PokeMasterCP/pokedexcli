package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

	return areas, nil
}
