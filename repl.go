package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	var cleanOutput []string
	trimmed := strings.TrimSpace(text)
	words := strings.SplitSeq(trimmed, " ")
	for w := range words {
		if w == "" {
			continue
		}
		word := strings.ToLower(string(w))
		cleanOutput = append(cleanOutput, word)
	}
	return cleanOutput
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		cmdName := input[0]
		if cmd, ok := commands[cmdName]; ok {
			cmd.callback()
		} else {
			fmt.Println("unknown command")
		}

	}
}

type cliCmd struct {
	name        string
	description string
	callback    func() error
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
	}
}
