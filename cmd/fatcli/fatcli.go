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
		ExitCommand{},
	}
	p := prompt.New(executor, completer)
	p.Run()
}

func completer(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "status set train {seconds}", Description: "Sets the status to training with the given {seconds} to train for"},
		{Text: "status set stop", Description: "Stops the training"},
		{Text: "status get", Description: "Returns the current status"},
	}
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
