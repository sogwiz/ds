package actions

import (
	"ds/pkg/master"
	"ds/pkg/slave"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
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

func GetFileFromDB(c *cli.Context) error {
	panic("implement me")
}

func PutFileInDB(c *cli.Context) error {
	host := c.String("host")
	port := c.Int("port")
	masterHostname := host + ":" + strconv.Itoa(port)

	path := c.String("file")
	if path == "" {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, "go", "src", "ds", "cmd", "ds", "main.go")
	}

	f, err := os.Open(path)
	if err != nil {
		logrus.Fatal("file not found", err)
	}
	defer f.Close()

	conn, err := net.Dial("tcp", masterHostname)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, fname := filepath.Split(path)

	_, _ = conn.Write([]byte("user_1/" + fname + "|"))
	_, _ = io.Copy(conn, f)
}
