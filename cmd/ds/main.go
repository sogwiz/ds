package main

import (
	"ds/pkg/actions"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	logrus.SetReportCaller(true)
	app := &cli.App{
		HideHelp:    true,
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Usage: "master node host",
				Value: "127.0.0.1",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "master node port",
				Value:   3300,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "display debug logs",
			},
			&cli.BoolFlag{
				Name:  "help",
				Usage: "show help",
			},
			&cli.BoolFlag{
				Name:  "version",
				Usage: "print the version",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("help") {
				cli.ShowAppHelpAndExit(c, 0)
			}
			if c.Bool("version") {
				cli.ShowVersion(c)
				os.Exit(0)
			}
			if c.Bool("verbose") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
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
			{
				Name:   "get",
				Usage:  "Get a file from the database",
				Action: actions.GetFileFromDB,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "file to get from the database",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "if defined, will store the file at this location (otherwise stdout)",
					},
				},
			},
			{
				Name:   "put",
				Usage:  "Put a file in the database",
				Action: actions.PutFileInDB,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "file to put in the database",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
