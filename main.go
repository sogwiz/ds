package main

import (
	"ds/pkg/actions"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"os"
)


func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			fmt.Println("run master or slave")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "master",
				Usage:   "Starts the master node",
				Action: actions.StartMaster,
			},
			{
				Name:    "slave",
				Usage:   "Starts the slave node",
				Action:  actions.StartSlave,
			},
		},
	}

	app.Run(os.Args)
}
