package master

import (
	"bufio"
	"ds/pkg/master/metadata"
	"ds/pkg/utils"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var meta *metadata.Metadata

func init() {
	meta = metadata.NewMetadata()
	meta.SetNumReplica(3)
	meta.AddNodes([]metadata.HostName{"slave1:3333", "slave2:3333", "slave3:3333", "slave4:3333", "slave5:3333", "slave6:3333"})
	meta.SetFileNodes(metadata.FileName("user_1/myfile.txt"), []metadata.HostName{"slave1:3333"})
}

func GetFile(filename metadata.FileName, conn net.Conn) {
	logrus.Debugf("GetFile %s", filename)
	nodes, exists := meta.GetFileNodes(filename)
	if !exists {
		_, _ = conn.Write([]byte("file not found\n"))
		conn.Close()
		return
	}
	randomNode := nodes.Random()
	logrus.Debugf("get file from random node %s", randomNode)
	nodeConn, err := net.Dial("tcp", string(randomNode))
	if err != nil {
		panic(err)
	}
	defer nodeConn.Close()
	_, _ = nodeConn.Write([]byte("GET|" + string(filename) + "|"))
	_, _ = io.Copy(conn, nodeConn)
	conn.Close()
}

func PutFile(filename metadata.FileName, fileContentStream io.Reader) {
	fileNodes := meta.GetOrCreateFileNodes(filename)
	firstNodeHostname := fileNodes.Shift()

	// Open connection to node1
	conn, err := net.Dial("tcp", string(firstNodeHostname))
	if err != nil {
		panic(err)
	}

	// TODO: could use some stream compression or blocks compression (lz4 ?)
	_, _ = conn.Write([]byte("PUT|"))
	_, _ = conn.Write([]byte(string(filename) + "|"))
	_, _ = conn.Write([]byte(fileNodes.Encode() + "|"))
	_, _ = io.Copy(conn, fileContentStream)

	fmt.Println("nodes:", fileNodes)
}

func CreateNewSlaveNode(ip string) {
	meta.AddNode(metadata.HostName(ip))
}

func handleRequest(conn net.Conn) {
	logrus.Debug("master handle conn")
	reader := bufio.NewReader(conn)
	method, _ := reader.ReadString('|')
	method = strings.TrimSuffix(method, "|")
	filename, _ := reader.ReadString('|')
	filename = strings.TrimSuffix(filename, "|")
	if method == "GET" {
		GetFile(metadata.FileName(filename), conn)
	} else if method == "PUT" {
		PutFile(metadata.FileName(filename), reader)
	}
}

// StartTCPServer ...
func StartTCPServer(host string, port int) {
	ctx := utils.SignalCtx()
	logrus.Info("master node listening on " + host + ":" + strconv.Itoa(port))
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
