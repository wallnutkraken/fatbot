package main

import (
	"os"
	"strconv"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/wallnutkraken/fatbot/cmd/fatcli/fatcaller"
	"github.com/wallnutkraken/fatbot/fatctrl/ctrltypes"
)

func main() {
	cl := fatcaller.New("http://www.fatbot.cli:1587")
	args := os.Args
	switch args[1] {
	case "status":
		switch args[2] {
		case "set":
			if args[3] == "train" {
				seconds, err := strconv.Atoi(args[4])
				if err != nil {
					logrus.WithError(err).Fatal("Failed reading train seconds")
				}
				if err := cl.StartTraining(ctrltypes.StartTrainingRequest{
					EndAfterSeconds: seconds,
				}); err != nil {
					logrus.WithError(err).Fatal("Request failed")
				}
			} else if args[3] == "stop" {
				if err := cl.StopTraining(); err != nil {
					logrus.WithError(err).Fatal("Request failed")
				}
			}
		case "get":
			status, err := cl.GetStatus()
			if err != nil {
				logrus.WithError(err).Fatal("Request failed")
			}
			fmt.Printf("Neural network status: %s\n", status.Network)
		}
	}
}
