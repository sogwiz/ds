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
	"strings"

	"github.com/sirupsen/logrus"
)

var meta *metadata.Metadata

func init() {
	meta = metadata.NewMetadata()
	meta.SetNumReplica(3)
	//meta.AddNodes([]metadata.HostName{"slave1:3333", "slave2:3334", "slave3:3335", "slave4:3336", "slave5:3337", "slave6:3338"})
	meta.AddNodes([]metadata.HostName{"slave1:3333", "slave2:3333", "slave3:3333", "slave4:3333", "slave5:3333", "slave6:3333"})
	//meta.AddNodes([]metadata.HostName{"slave1", "slave2", "slave3", "slave4", "slave5", "slave6"})
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
	// Read 1024 bytes at the time and stream it to slave
	buf := make([]byte, 1024)
	for {
		n, err := fileContentStream.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		if _, err := conn.Write(buf[:n]); err != nil {
			logrus.Error(err)
		}
	}

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
func StartTCPServer() {
	c1, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})

	go func(ctx context.Context) {
		logrus.Info("slave server listening on 0.0.0.0:3333")
		l, err := net.Listen("tcp", "0.0.0.0:3333")
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
