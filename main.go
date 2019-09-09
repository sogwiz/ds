package main

import (
	"ds/pkg/actions"
	"gopkg.in/urfave/cli.v2"
	"fmt"
	"os"
)


func main() {
	app := &cli.App{
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "lang",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
		Action: func(c *cli.Context) error {
			name := "Nefertiti"
			if c.NArg() > 0 {
				name = c.Args().Get(0)
			}
			if c.String("lang") == "spanish" {
				fmt.Println("Hola", name)
			} else {
				fmt.Println("Hello", name)
			}
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
