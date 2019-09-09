package actions

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
)

func StartMaster(c *cli.Context) error {
	 fmt.Println("Starts master")
	return nil
}

func StartSlave(c *cli.Context) error {
	 fmt.Println("Starts slave")
	return nil
}
