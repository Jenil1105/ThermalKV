package index

import (
	"sync"
)

type ColdEntry struct {
	Offset int64
	Expiry int64
}

type ColdIndex struct {
	ColdIndex map[string]ColdEntry
	Mutex     sync.RWMutex
}

func NewColdIndex() *ColdIndex {

	return &ColdIndex{
		ColdIndex: make(map[string]ColdEntry),
	}

}

func (i *ColdIndex) AddColdIndex(key string, offset, expiry int64) {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	i.ColdIndex[key] = ColdEntry{
		Offset: offset,
		Expiry: expiry,
	}
}
