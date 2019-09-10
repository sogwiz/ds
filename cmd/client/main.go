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
		Action: func(c *cli.Context) error {
			logrus.Info("ds client")
			return nil
		},
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
			return nil
		},
		Commands: []*cli.Command{
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
