package main

import "strings"

func splitOmitEmpty(text string) []string {
	args := strings.Split(text, " ")
	fullArgs := []string{}
	for _, arg := range args {
		if strings.TrimSpace(arg) != "" {
			fullArgs = append(fullArgs, arg)
		}
	}

	return fullArgs
}

type Command interface {
	Is(text string) bool
	Exec(text string) error
}
