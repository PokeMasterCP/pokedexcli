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
		var param string
		if len(input) == 0 {
			continue
		} else if len(input) > 1 {
			param = input[1]
		}

		cmdName := input[0]
		if cmd, ok := commands[cmdName]; ok {
			err := cmd.callback(&config, param)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("unknown command")
		}

	}
}
