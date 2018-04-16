package main

import "os"

type ExitCommand struct{}

func (ExitCommand) Is(text string) bool {
	return text == "exit" || text == "quit"
}

func (ExitCommand) Exec(text string) error {
	os.Exit(0)
	return nil
}
