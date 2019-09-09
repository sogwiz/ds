package master

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

type UserID int64

type HostName string

type Metadata struct {
	sync.RWMutex
	numReplica  int32 // atomic value
	users       map[UserID][]HostName
	allNodesMap map[HostName]bool
}

func NewMetadata() *Metadata {
	return new(Metadata)
}

func (m *Metadata) GetNumReplica() int32 {
	return atomic.LoadInt32(&m.numReplica)
}

func (m *Metadata) SetNumReplica(num int32) {
	atomic.StoreInt32(&m.numReplica, num)
}

func (m *Metadata) AddNode(hostname HostName) {
	m.Lock()
	defer m.Unlock()
	m.allNodesMap[hostname] = true
}

func (m *Metadata) AddNodes(hostnames []HostName) {
	m.Lock()
	defer m.Unlock()
	for _, hostname := range hostnames {
		m.allNodesMap[hostname] = true
	}
}

func (m *Metadata) GetNodesCount() int32 {
	m.Lock()
	defer m.Unlock()
	return int32(len(m.allNodesMap))
}

func (m *Metadata) SetUserNodes(userID UserID, nodes []HostName) {
	m.Lock()
	defer m.Unlock()
	m.users[userID] = nodes
}

func (m *Metadata) GetUserNodes(userID UserID) (hostnames []HostName, exists bool) {
	m.Lock()
	defer m.Unlock()
	hostnames, exists = m.users[userID]
	return
}

// Pick a random hostname from the meta nodes map
func (m *Metadata) GetRandomHostName() HostName {
	m.Lock()
	defer m.Unlock()
	i := rand.Intn(len(meta.allNodesMap))
	for k := range meta.allNodesMap {
		if i == 0 {
			return k
		}
		i--
	}
	panic("never")
}

var meta *Metadata

func init() {
	meta = NewMetadata()
	meta.SetNumReplica(3)
	meta.AddNodes([]HostName{"IP1", "IP2", "IP3", "IP4", "IP5", "IP6"})
	meta.SetUserNodes(UserID(1), []HostName{"IP1", "IP2", "IP3"})
}

// Generates "num" random unique indexes
func generateRandomHostnames(num int32) (hostnames []HostName) {
	hostnamesMap := make(map[HostName]bool)
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

func createUserInMetadata(userID UserID) {
	if meta.GetNodesCount() < meta.GetNumReplica() {
		panic("not enough nodes, need at least 3")
	}
	meta.SetUserNodes(userID, generateRandomHostnames(meta.GetNumReplica()))
}

func PutFile(userID UserID, fileName string, fileContentStream io.Reader) {
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
	meta.AddNode(HostName(ip))
}
