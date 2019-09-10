package config

import "sync"

var inst *Config

// GetInstance ...
func GetInstance() *Config {
	if inst == nil {
		inst = new(Config)
		return inst
	}
	return inst
}

// Config ...
type Config struct {
	sync.RWMutex
	dataPath string
}

// GetDataPath ...
func (c *Config) GetDataPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.dataPath
}

// SetDataPath ...
func (c *Config) SetDataPath(newPath string) {
	c.Lock()
	defer c.Unlock()
	c.dataPath = newPath
}
