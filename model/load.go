package model

import "sync"

type LoadMap map[Ticker][]*Period

var (
	lock = sync.RWMutex{}
)

func (m LoadMap) Read(key Ticker) []*Period {
	lock.RLock()
	defer lock.RUnlock()
	return m[key]
}

func (m LoadMap) Write(key Ticker, value []*Period) {
	lock.Lock()
	defer lock.Unlock()
	m[key] = value
}
