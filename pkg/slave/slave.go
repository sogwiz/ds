package slave

import (
	"bufio"
	"ds/pkg/config"
	"ds/pkg/utils"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func handleGetRequest(conn net.Conn, reader *bufio.Reader) {
	filename, _ := reader.ReadString('|')
	filename = strings.TrimSuffix(filename, "|")

	path := filepath.Join(config.GetInstance().GetDataPath(), filename)
	fo, err := os.Open(path)
	if err != nil {
		_, _ = conn.Write([]byte("file not found\n"))
		return
	}
	defer fo.Close()
	_, _ = io.Copy(conn, fo)
}

func handlePutRequest(conn net.Conn, reader *bufio.Reader) {
	filename, _ := reader.ReadString('|')
	filename = strings.TrimSuffix(filename, "|")

	hostnamesRaw, _ := reader.ReadString('|')
	hostnamesRaw = strings.TrimSuffix(hostnamesRaw, "|")

	hostnames := strings.Split(hostnamesRaw, ",")

	dir, file := filepath.Split(filename)
	dataPath := config.GetInstance().GetDataPath()
	_ = os.MkdirAll(filepath.Join(dataPath, dir), 0777)
	fo, err := os.Create(filepath.Join(dataPath, dir, file))
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if len(hostnames) == 0 || hostnames[0] == "" {
		_, _ = io.Copy(fo, reader)
	} else {
		copyToHostname := hostnames[0]                       // We need to copy the file to this guy
		hostnamesEncoded := strings.Join(hostnames[1:], ",") // This is the list of hostnames to pass along

		nextConn, err := net.Dial("tcp", copyToHostname)
		if err != nil {
			panic(err)
		}
		defer nextConn.Close()
		_, _ = nextConn.Write([]byte("PUT|"))
		_, _ = nextConn.Write([]byte(filename + "|"))
		_, _ = nextConn.Write([]byte(hostnamesEncoded + "|"))
		if _, err := io.Copy(io.MultiWriter(fo, nextConn), reader); err != nil {
			panic(err)
		}
	}

	fmt.Println("received a file: " + filename)
}

func handleRequest(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	reader := bufio.NewReader(conn)

	method, _ := reader.ReadString('|')
	method = strings.TrimSuffix(method, "|")

	switch method {
	case "GET":
		handleGetRequest(conn, reader)
	case "PUT":
		handlePutRequest(conn, reader)
	}
}

func StartTCPServer(host string, port int) {
	ctx := utils.SignalCtx()
	logrus.Info("slave node listening on " + host + ":" + strconv.Itoa(port))
	l, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			logrus.Info("cancelled")
			return
		case conn := <-utils.AcceptConn(l):
			go handleRequest(conn)
		}
	}
}
