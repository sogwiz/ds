package master

import (
	"bufio"
	"context"
	"ds/pkg/master/metadata"
	"ds/pkg/utils"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var meta *metadata.Metadata

func init() {
	meta = metadata.NewMetadata()
	meta.SetNumReplica(3)
	meta.AddNodes([]metadata.HostName{"slave1:3333", "slave2:3333", "slave3:3333", "slave4:3333", "slave5:3333", "slave6:3333"})
}

func PutFile(filename metadata.FileName, fileContentStream io.Reader) {
	fileNodes := meta.GetOrCreateFileNodes(filename)

	// TODO: linked list, so datanode transfer file with each others instead of master node
	// Open connection to node1
	fmt.Println(fileNodes)
	conn, err := net.Dial("tcp", string(fileNodes[0]))
	if err != nil {
		panic(err)
	}

	_, _ = conn.Write([]byte(string(filename) + "|"))

	remainingHostnamesStr := fileNodes[1:]
	hostnamesStr := make([]string, 0)
	for _, hn := range remainingHostnamesStr {
		hostnamesStr = append(hostnamesStr, string(hn))
	}
	hostnamesEncoded := strings.Join(hostnamesStr, ",")
	_, _ = conn.Write([]byte(hostnamesEncoded + "|"))

	// TODO: could use some stream compression or blocks compression (lz4 ?)
	_, _ = io.Copy(conn, fileContentStream)

	// Stream file content to it with metadata about replicas

	fmt.Println("nodes:", fileNodes)
}

func CreateNewSlaveNode(ip string) {
	meta.AddNode(metadata.HostName(ip))
}

func handleRequest(conn net.Conn) {
	fmt.Println("Handle request")
	reader := bufio.NewReader(conn)
	filename, _ := reader.ReadString('|')
	filename = strings.TrimSuffix(filename, "|")

	PutFile(metadata.FileName(filename), reader)
}

// StartTCPServer ...
func StartTCPServer(host string, port int) {
	c1, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})

	go func(ctx context.Context) {
		logrus.Info("master node listening on " + host + ":" + strconv.Itoa(port))
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
