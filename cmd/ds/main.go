package main

import (
	"ds/pkg/actions"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			logrus.Info("run master or slave")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "master",
				Usage:  "Starts the master node",
				Action: actions.StartMaster,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Master host",
						Value: "127.0.0.1",
					},
					&cli.IntFlag{
						Name:  "port",
						Usage: "Master port",
						Value: 3300,
					},
				},
			},
			{
				Name:   "slave",
				Usage:  "Starts the slave node",
				Action: actions.StartSlave,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Slave host",
						Value: "127.0.0.1",
					},
					&cli.IntFlag{
						Name:  "port",
						Usage: "Slave port",
						Value: 3333,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
