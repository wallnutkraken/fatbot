package main

import (
	"flag"

	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/wallnutkraken/fatbot/cmd/fatcli/fatcaller"
)

var commands []Command

func main() {
	hostname := flag.String("h", "http://localhost:1587", "Hostname of the fatbot service")
	flag.Parse()
	cl := fatcaller.New(*hostname)
	commands = []Command{
		NewStatus(cl),
	}
	p := prompt.New(executor, completer)
	p.Run()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "status", Description: "Gets or sets the status of the neural network"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(cmd string) {
	for _, command := range commands {
		if command.Is(cmd) {
			if err := command.Exec(cmd); err != nil {
				fmt.Println("Error running command:", err.Error())
			} else {
				return
			}
		}
	}

	fmt.Println("Command not found")
}
