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
		word := strings.ToLower(w)
		cleanOutput = append(cleanOutput, word)
	}
	return cleanOutput
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	config := initCommands()

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		cmdName := input[0]
		if cmd, ok := commands[cmdName]; ok {
			err := cmd.callback(&config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("unknown command")
		}

	}
}
