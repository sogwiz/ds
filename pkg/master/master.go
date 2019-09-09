package master

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"

	"github.com/sirupsen/logrus"
)

type UserID int64

type HostName string

type Metadata struct {
	NumReplica  int
	Users       map[UserID][]HostName
	AllNodesMap map[HostName]bool
}

var meta Metadata

func init() {
	meta = Metadata{}
	meta.NumReplica = 3
	meta.AllNodesMap = map[HostName]bool{
		"IP1": true,
		"IP2": true,
		"IP3": true,
		"IP4": true,
		"IP5": true,
		"IP6": true,
	}
	meta.Users = make(map[UserID][]HostName)
	meta.Users[UserID(1)] = []HostName{"IP1", "IP2", "IP3"}
}

// Pick a random hostname from the meta nodes map
func randHostName() HostName {
	i := rand.Intn(len(meta.AllNodesMap))
	for k := range meta.AllNodesMap {
		if i == 0 {
			return k
		}
		i--
	}
	panic("never")
}

// Generates "num" random unique indexes
func generateRandomHostnames(num int) (indexesArr []HostName) {
	indexesMap := make(map[HostName]bool)
	for len(indexesMap) != num {
		tmpHostName := randHostName()
		_, exists := indexesMap[tmpHostName]
		if exists {
			continue
		}
		indexesArr = append(indexesArr, tmpHostName)
		indexesMap[tmpHostName] = true
	}
	return
}

func createUserInMetadata(userID UserID) {
	if len(meta.AllNodesMap) < meta.NumReplica {
		panic("not enough nodes, need at least 3")
	}
	indexesArray := generateRandomHostnames(meta.NumReplica)
	// Add Nodes to user nodes
	meta.Users[userID] = make([]HostName, 0)
	for _, hostname := range indexesArr {
		meta.Users[userID] = append(meta.Users[userID], hostname)
	}
}

func PutFile(userID UserID, fileName string, fileContentStream io.Reader) {
	userNodes, userExists := meta.Users[userID]
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
	meta.AllNodesMap[HostName(ip)] = true
}
