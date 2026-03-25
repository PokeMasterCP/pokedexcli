package main

import (
	"fmt"
	"os"

	"github.com/pokemastercp/pokedexcli/internal/pokeapi"
)

type config struct {
	next *string
	prev *string
}

func initCommands() config {
	var conf config

	locationAreaStart := "https://pokeapi.co/api/v2/location-area?limit=20&offset=0"
	conf.next = &locationAreaStart
	return conf
}

type cliCmd struct {
	name        string
	description string
	callback    func(c *config) error
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
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, command := range commands {
		fmt.Printf("\n%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(c *config) error {
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

func commandMapb(c *config) error {
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
