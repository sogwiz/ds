package actions

import (
	"ds/pkg/master"
	"ds/pkg/slave"

	"gopkg.in/urfave/cli.v2"
)

func StartMaster(c *cli.Context) error {
	master.StartTCPServer()
	return nil
}

func StartSlave(c *cli.Context) error {
	slave.StartTCPServer()
	return nil
}
