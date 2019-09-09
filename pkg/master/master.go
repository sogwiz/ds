package master

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type UserID int64

type HostName string

type Metadata struct {
	sync.RWMutex
	numReplica  int
	users       map[UserID][]HostName
	allNodesMap map[HostName]bool
}

func (m *Metadata) AddNode(hostname HostName) {
	m.Lock()
	defer m.Unlock()
	m.allNodesMap[hostname] = true
}

var meta Metadata

func init() {
	meta = Metadata{}
	meta.numReplica = 3
	meta.allNodesMap = map[HostName]bool{
		"IP1": true,
		"IP2": true,
		"IP3": true,
		"IP4": true,
		"IP5": true,
		"IP6": true,
	}
	meta.users = make(map[UserID][]HostName)
	meta.users[UserID(1)] = []HostName{"IP1", "IP2", "IP3"}
}

// Pick a random hostname from the meta nodes map
func randHostName() HostName {
	i := rand.Intn(len(meta.allNodesMap))
	for k := range meta.allNodesMap {
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
	if len(meta.allNodesMap) < meta.numReplica {
		panic("not enough nodes, need at least 3")
	}
	meta.users[userID] = generateRandomHostnames(meta.numReplica)
}

func PutFile(userID UserID, fileName string, fileContentStream io.Reader) {
	userNodes, userExists := meta.users[userID]
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
	meta.AddNode(HostName(ip))
}
