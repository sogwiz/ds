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
	meta.AddNodes([]metadata.HostName{"localhost:3333", "localhost:3334", "localhost:3335", "localhost:3336", "localhost:3337", "localhost:3338"})
}

func PutFile(filename metadata.FileName, fileContentStream io.Reader) {
	fileNodes := meta.GetOrCreateFileNodes(filename)

	for _, fileNode := range fileNodes {
		// Open connection to node1
		conn, err := net.Dial("tcp", string(fileNode))
		if err != nil {
			panic(err)
		}

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

		fmt.Println("nodes:", fileNodes)
	}
}

func CreateNewSlaveNode(ip string) {
	meta.AddNode(metadata.HostName(ip))
}
