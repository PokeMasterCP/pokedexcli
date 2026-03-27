package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/pokemastercp/pokedexcli/internal/pokeapi"
)

type config struct {
	next    *string
	prev    *string
	pokedex map[string]pokeapi.Pokemon
}

func initCommands() config {
	var conf config

	locationAreaStart := "https://pokeapi.co/api/v2/location-area?limit=20&offset=0"
	conf.next = &locationAreaStart
	conf.pokedex = make(map[string]pokeapi.Pokemon)
	return conf
}

type cliCmd struct {
	name        string
	description string
	callback    func(c *config, param string) error
}

func getCommands() map[string]cliCmd {
	return map[string]cliCmd{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Print next page of maps",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Print previous page of maps",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area's Pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon's info",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List your caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(c *config, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, param string) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, command := range commands {
		fmt.Printf("\n%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(c *config, param string) error {
	if c.next == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	url := *c.next
	areas, err := pokeapi.GetLocationArea(url)
	if err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	c.next = areas.Next
	c.prev = areas.Previous
	return nil
}

func commandMapb(c *config, param string) error {
	if c.prev == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := *c.prev
	areas, err := pokeapi.GetLocationArea(url)
	if err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	c.next = areas.Next
	c.prev = areas.Previous
	return nil
}

func commandExplore(c *config, param string) error {
	if param == "" {
		return fmt.Errorf("provide an area to explore")
	}

	locationPokemon, err := pokeapi.GetLocationPokemon(param)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %v...\nFound Pokemon:\n", param)
	for _, entry := range locationPokemon.PokemonEncounters {
		fmt.Printf("- %v\n", entry.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config, param string) error {
	if param == "" {
		return fmt.Errorf("provide a pokemon name to catch")
	}
	name := strings.ToLower(param)

	_, caught := c.pokedex[name]
	if caught {
		fmt.Printf("You've already caught %v!\n", name)
		return nil
	}

	pokemonData, err := pokeapi.GetPokemonData(name)
	if err != nil {
		return err
	}

	baseExp := pokemonData.BaseExperience
	fmt.Printf("Throwing a Pokeball at %v...\n", name)

	successfulCatch := baseExp < rand.Intn(650)
	if successfulCatch {
		fmt.Printf("%v was caught!\n", name)
		c.pokedex[name] = pokemonData
	} else {
		fmt.Printf("%v escaped!\n", name)
	}

	return nil
}

func commandInspect(c *config, param string) error {
	if param == "" {
		return fmt.Errorf("return a pokemon name to inspect")
	}
	name := strings.ToLower(param)

	data, caught := c.pokedex[name]
	if !caught {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %v\n", data.Name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Println("Stats:")

	for _, stat := range data.Stats {
		fmt.Printf("  -%v: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range data.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *config, param string) error {
	if len(c.pokedex) == 0 {
		return fmt.Errorf("you haven't caught any Pokemon yet!")
	}

	fmt.Println("Your Pokedex:")
	for key := range c.pokedex {
		fmt.Printf("  -%s\n", key)
	}

	return nil
}
