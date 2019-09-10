package slave

import (
	"bufio"
	"context"
	"ds/pkg/utils"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	filename, _ := reader.ReadString('|')
	filename = strings.TrimSuffix(filename, "|")

	hostnamesRaw, _ := reader.ReadString('|')
	hostnamesRaw = strings.TrimSuffix(hostnamesRaw, "|")

	hostnames := strings.Split(hostnamesRaw, ",")

	needToCopy := false

	var hostnamesEncoded string
	var nextConn net.Conn

	if len(hostnames) == 0 || hostnames[0] == "" {
		logrus.Debug("should not copy to anyone else")
	} else {
		needToCopy = true

		copyToHostname := hostnames[0]                      // We need to copy the file to this guy
		hostnamesEncoded = strings.Join(hostnames[1:], ",") // This is the list of hostnames to pass along

		var err error
		nextConn, err = net.Dial("tcp", copyToHostname)
		if err != nil {
			panic(err)
		}
	}

	// Read the incoming connection into the buffer.
	homedir, _ := os.UserHomeDir()
	dir, file := filepath.Split(filename)
	_ = os.MkdirAll(filepath.Join(homedir, "data", dir), 0777)
	fo, err := os.Create(filepath.Join(homedir, "data", dir, file))
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if needToCopy {
		_, _ = nextConn.Write([]byte(filename + "|"))
		_, _ = nextConn.Write([]byte(hostnamesEncoded + "|"))
		if _, err := io.Copy(io.MultiWriter(fo, nextConn), reader); err != nil {
			panic(err)
		}
	} else {
		_, _ = io.Copy(fo, reader)
	}

	fmt.Println("received a file: " + filename)

	if err := conn.Close(); err != nil {
		logrus.Error(err)
	}
}

func StartTCPServer(host string, port int) {
	c1, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})

	go func(ctx context.Context) {
		logrus.Info("slave node listening on " + host + ":" + strconv.Itoa(port))
		l, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			panic(err)
		}
		for {
			select {
			case <-ctx.Done():
				logrus.Info("cancelled")
				close(exitCh)
				return
			case conn := <-utils.AcceptConn(l):
				go handleRequest(conn)
			}
		}
	}(c1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-exitCh
}
