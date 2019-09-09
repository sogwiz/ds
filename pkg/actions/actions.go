package actions

import (
	"ds/pkg/slave"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func StartMaster(c *cli.Context) error {
	logrus.Info("Starts master")
	return nil
}

func StartSlave(c *cli.Context) error {
	slave.StartTCPServer()
	return nil
}
