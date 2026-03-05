package store

import (
	"sync"

	"github.com/brandondvs/flick/internal/feature"
)

type Memory struct {
	data map[string]*feature.Flag

	sync.RWMutex
}

func New() *Memory {
	return &Memory{
		data: make(map[string]*feature.Flag),
	}
}

func (m *Memory) Store(name string, featureFlag *feature.Flag) {
	m.Lock()
	m.data[name] = featureFlag
	m.Unlock()
}

func (m *Memory) Get(name string) *feature.Flag {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.data[name]
	if !ok {
		return nil
	}
	return v
}

func (m *Memory) Delete(name string) {
	m.Lock()
	delete(m.data, name)
	m.Unlock()
}

func (m *Memory) AllKeys() []string {
	m.RLock()
	defer m.RUnlock()
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}
