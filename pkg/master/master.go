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

type Node struct {
	IP string
}

type Metadata struct {
	NumReplica int
	Users      map[UserID][]Node
	AllNodes   []Node
}

var meta Metadata

func init() {
	meta = Metadata{}
	meta.NumReplica = 3
	meta.AllNodes = []Node{{"IP1"}, {"IP2"}, {"IP3"}, {"IP4"}, {"IP5"}, {"IP6"}}
	meta.Users = make(map[UserID][]Node)
	meta.Users[UserID(1)] = []Node{{"IP1"}, {"IP2"}, {"IP3"}}
}

// Generates "num" random unique indexes
func generateRandomIndexes(num int) (indexesArr []int) {
	indexesMap := make(map[int]bool)
	for len(indexesMap) != num {
		tmpNum := rand.Intn(len(meta.AllNodes))
		_, exists := indexesMap[tmpNum]
		if exists {
			continue
		}
		indexesArr = append(indexesArr, tmpNum)
		indexesMap[tmpNum] = true
	}
	return
}

func createUserInMetadata(userID UserID) {
	if len(meta.AllNodes) < meta.NumReplica {
		panic("not enough nodes, need at least 3")
	}
	indexesArr := generateRandomIndexes(meta.NumReplica)
	// Add Nodes to user nodes
	meta.Users[userID] = make([]Node, 0)
	for _, index := range indexesArr {
		meta.Users[userID] = append(meta.Users[userID], meta.AllNodes[index])
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
	// TODO: Add validation if node already exists. (should we use a map instead of an array)
	meta.AllNodes = append(meta.AllNodes, Node{IP: ip})
}
