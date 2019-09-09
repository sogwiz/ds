package metadata

import (
	"math/rand"
	"sync"
	"sync/atomic"
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
	i := rand.Intn(len(m.allNodesMap))
	for k := range m.allNodesMap {
		if i == 0 {
			return k
		}
		i--
	}
	panic("never")
}
