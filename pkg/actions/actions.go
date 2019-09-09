package actions

import (
	"ds/pkg/slave"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func StartMaster(c *cli.Context) error {
	logrus.Info("Starts master")
	for {
		time.Sleep(10 * time.Second)
	}
	return nil
}

func StartSlave(c *cli.Context) error {
	slave.StartTCPServer()
	return nil
}
