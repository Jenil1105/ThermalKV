package store

import (
	"container/heap"
	"sync"
	"thermalkv/internal/coldstore"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/ttl"
)

type Store struct {
	Data               map[string]model.Item
	CurrentMemoryUsage int64
	MaxHotMemory       int64
	CoolingThreshold   int64
	Mutex              sync.RWMutex
	ExpiryHeap         ttl.MinHeap
	WAL                *walpkg.WAL
	Thermal            *coldstore.Manager
	CoolingInProgress  bool
}

// NewStore
func NewStore(wal *walpkg.WAL, manager *coldstore.Manager) *Store {
	h := ttl.MinHeap{}
	heap.Init(&h)
	return &Store{
		Data:              make(map[string]model.Item),
		MaxHotMemory:      100,
		CoolingThreshold:  100000,
		ExpiryHeap:        h,
		WAL:               wal,
		Thermal:           manager,
		CoolingInProgress: false,
	}
}
