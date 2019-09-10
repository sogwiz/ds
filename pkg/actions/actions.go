package actions

import (
	"ds/pkg/master"
	"ds/pkg/slave"
	"ds/pkg/utils"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// BeforeApp this hook is called before any other actions
func BeforeApp(c *cli.Context) error {
	logrus.SetFormatter(utils.LogFormatter{})
	if c.Bool("help") {
		cli.ShowAppHelpAndExit(c, 0)
	}
	if c.Bool("version") {
		cli.ShowVersion(c)
		os.Exit(0)
	}
	if c.Bool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("verbose mode enabled")
	}
	return nil
}

func StartMaster(c *cli.Context) error {
	host := c.String("host")
	port := c.Int("port")
	master.StartTCPServer(host, port)
	return nil
}

func StartSlave(c *cli.Context) error {
	host := c.String("host")
	port := c.Int("port")
	slave.StartTCPServer(host, port)
	return nil
}

func GetFileFromDB(c *cli.Context) error {
	host := c.String("host")
	port := c.Int("port")
	masterHostname := host + ":" + strconv.Itoa(port)
	filePath := c.String("file")
	o := os.Stdout
	output := c.String("output")
	var fo *os.File
	if output != "" {
		var err error
		fo, err = os.Create(output)
		if err != nil {
			logrus.Fatal("failed to create output file: ", err)
		}
		defer fo.Close()
		o = fo
	}
	conn, err := net.Dial("tcp", masterHostname)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, _ = conn.Write([]byte("GET|user_1/" + filePath + "|"))
	_, _ = io.Copy(o, conn)
	return nil
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

	_, _ = conn.Write([]byte("PUT|user_1/" + fname + "|"))
	_, _ = io.Copy(conn, f)
	return nil
}
