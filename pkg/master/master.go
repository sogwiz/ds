package master

import (
	"ds/pkg/master/metadata"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

var meta *metadata.Metadata

func init() {
	meta = metadata.NewMetadata()
	meta.SetNumReplica(3)
	meta.AddNodes([]metadata.HostName{"IP1", "IP2", "IP3", "IP4", "IP5", "IP6"})
}

// Generates "num" random unique indexes
func generateRandomHostnames(num int32) (hostnames []metadata.HostName) {
	hostnamesMap := make(map[metadata.HostName]bool)
	for int32(len(hostnamesMap)) != num {
		tmpHostName := meta.GetRandomHostName()
		_, exists := hostnamesMap[tmpHostName]
		if exists {
			continue
		}
		hostnames = append(hostnames, tmpHostName)
		hostnamesMap[tmpHostName] = true
	}
	return
}

func createfileInMetadata(file metadata.FileName) {
	if meta.GetNodesCount() < meta.GetNumReplica() {
		panic("not enough nodes, need at least 3")
	}
	meta.SetFileNodes(file, generateRandomHostnames(meta.GetNumReplica()))
}

func PutFile(filename metadata.FileName, fileContentStream io.Reader) {
	userNodes, userExists := meta.GetFileNodes(filename)
	if !userExists {
		createfileInMetadata(filename)
	}

	// Open connection to node1
	conn, _ := net.Dial("tcp", "localhost:3333")

	_, _ = conn.Write([]byte(string(filename) + "|"))

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

	fmt.Println(userNodes)
}

func CreateNewSlaveNode(ip string) {
	meta.AddNode(metadata.HostName(ip))
}
