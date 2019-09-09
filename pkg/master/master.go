package master

import (
	"ds/pkg/master/metadata"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/sirupsen/logrus"
)

var meta *metadata.Metadata

func init() {
	meta = metadata.NewMetadata()
	meta.SetNumReplica(3)
	meta.AddNodes([]metadata.HostName{"IP1", "IP2", "IP3", "IP4", "IP5", "IP6"})
	meta.SetUserNodes(metadata.UserID(1), []metadata.HostName{"IP1", "IP2", "IP3"})
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

func createUserInMetadata(userID metadata.UserID) {
	if meta.GetNodesCount() < meta.GetNumReplica() {
		panic("not enough nodes, need at least 3")
	}
	meta.SetUserNodes(userID, generateRandomHostnames(meta.GetNumReplica()))
}

func PutFile(userID metadata.UserID, fileName string, fileContentStream io.Reader) {
	userNodes, userExists := meta.GetUserNodes(userID)
	if !userExists {
		createUserInMetadata(userID)
	}

	// Open connection to node1
	conn, _ := net.Dial("tcp", "localhost:3333")

	content, _ := ioutil.ReadAll(fileContentStream)
	_, err := conn.Write(content)
	if err != nil {
		logrus.Error(err)
	}

	// Stream file content to it with metadata about replicas

	fmt.Println(userNodes)
}

func CreateNewSlaveNode(ip string) {
	meta.AddNode(metadata.HostName(ip))
}
