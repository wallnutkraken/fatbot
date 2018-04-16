package main

import (
	"errors"
	"strings"

	"fmt"

	"strconv"

	"github.com/wallnutkraken/fatbot/cmd/fatcli/fatcaller"
	"github.com/wallnutkraken/fatbot/fatctrl/ctrltypes"
)

type StatusCommand struct {
	cl *fatcaller.Client
}

func NewStatus(cl *fatcaller.Client) *StatusCommand {
	return &StatusCommand{
		cl: cl,
	}
}

func (StatusCommand) Is(text string) bool {
	return strings.HasPrefix(text, "status")
}

func (s *StatusCommand) Exec(text string) error {
	args := splitOmitEmpty(text)
	if len(args) < 2 {
		return errors.New("Status command requires at least two arguments")
	}
	switch args[1] {
	case "get":
		return s.get(args)
	case "set":
		return s.set(args)
	default:
		return errors.New("Unrecognized status command argument: " + args[1])
	}
}

func (s *StatusCommand) set(args []string) error {
	if len(args) < 3 {
		return errors.New("Status command \"set\" requires at least three arguments")
	}
	switch args[2] {
	case "train":
		if len(args) < 4 {
			return errors.New("Status command \"set train\" requires a train-for argument (in seconds)")
		}
		secs, err := strconv.Atoi(args[3])
		if err != nil {
			return errors.New("train-for argument MUST be a number")
		}
		return s.cl.StartTraining(ctrltypes.StartTrainingRequest{EndAfterSeconds: secs})
	case "stop":
		return s.cl.StopTraining()
	default:
		return errors.New("Argument " + args[2] + " was not recognized")
	}
}

func (s *StatusCommand) get(args []string) error {
	status, err := s.cl.GetStatus()
	if err != nil {
		return err
	}
	fmt.Println("Neural network status:", status.Network)
	return nil
}
