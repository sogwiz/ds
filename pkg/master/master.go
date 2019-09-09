package master

import (
	"ds/pkg/master/metadata"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

var meta *metadata.Metadata

func init() {
	meta = metadata.NewMetadata()
	meta.SetNumReplica(3)
	meta.AddNodes([]metadata.HostName{"slave1:3333", "slave2:3334", "slave3:3335", "slave4:3336", "slave5:3337", "slave6:3338"})
}

func PutFile(filename metadata.FileName, fileContentStream io.Reader) {
	fileNodes := meta.GetOrCreateFileNodes(filename)

	// TODO: linked list, so datanode transfer file with each others instead of master node
	// Open connection to node1
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
