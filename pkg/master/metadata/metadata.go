package metadata

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

type UserID int64

type FileName string

type HostName string

type Metadata struct {
	sync.RWMutex
	numReplica  int32 // atomic value
	files       map[FileName][]HostName
	allNodesMap map[HostName]bool
}

func NewMetadata() *Metadata {
	m := new(Metadata)
	m.files = make(map[FileName][]HostName)
	m.allNodesMap = make(map[HostName]bool)
	return m
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

func (m *Metadata) SetFileNodes(file FileName, nodes []HostName) {
	m.Lock()
	defer m.Unlock()
	m.files[file] = nodes
}

func (m *Metadata) setFileNodes(file FileName, nodes []HostName) {
	m.files[file] = nodes
}

func (m *Metadata) GetOrCreateFileNodes(file FileName) (hostnames []HostName) {
	m.Lock()
	defer m.Unlock()
	hostnames, exists := m.getFileNodes(file)
	if exists {
		return hostnames
	}
	hostnames = m.generateRandomHostnames()
	m.setFileNodes(file, hostnames)
	return hostnames
}

func (m *Metadata) GetFileNodes(file FileName) (hostnames []HostName, exists bool) {
	m.Lock()
	defer m.Unlock()
	return m.getFileNodes(file)
}

func (m *Metadata) getFileNodes(file FileName) (hostnames []HostName, exists bool) {
	hostnames, exists = m.files[file]
	return
}

// Pick a random hostname from the meta nodes map
func (m *Metadata) getRandomHostName() HostName {
	i := rand.Intn(len(m.allNodesMap))
	for k := range m.allNodesMap {
		if i == 0 {
			return k
		}
		i--
	}
	panic("never")
}

// Generates "num" random unique indexes
func (m *Metadata) generateRandomHostnames() (hostnames []HostName) {
	hostnamesMap := make(map[HostName]bool)
	for int32(len(hostnamesMap)) != m.GetNumReplica() {
		tmpHostName := m.getRandomHostName()
		_, exists := hostnamesMap[tmpHostName]
		if exists {
			continue
		}
		hostnames = append(hostnames, tmpHostName)
		hostnamesMap[tmpHostName] = true
	}
	return
}
