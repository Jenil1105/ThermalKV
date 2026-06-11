package store

import (
	"container/heap"
	"sync"
	"thermalkv/internal/model"
	"thermalkv/internal/ttl"
)

type WALWriter interface {
	Write(operation string, key string, value ...string) error
}

type ColdStorage interface {
	MoveToCool(key string, item model.Item) error
	LoadFromCool(key string) (model.Item, bool)
	DeleteIndex(key string)
	HaveIndex(key string) bool
	GetIndexSize() int
	AppendDelete(key string) error
}

type Store struct {
	Data               map[string]model.Item
	CurrentMemoryUsage int64
	MaxHotMemory       int64
	CoolingThreshold   int64
	Mutex              sync.RWMutex
	ExpiryHeap         ttl.MinHeap
	WAL                WALWriter
	Thermal            ColdStorage
	CoolingInProgress  bool
}

// NewStore
func NewStore(wal WALWriter, manager ColdStorage) *Store {
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

func (s *Store) putItem(key string, item model.Item) {
	if oldItem, exists := s.Data[key]; exists {
		s.CurrentMemoryUsage -= oldItem.Size
	}

	s.Data[key] = item
	s.CurrentMemoryUsage += item.Size
}

func (s *Store) removeItem(key string) {
	if oldItem, exists := s.Data[key]; exists {
		s.CurrentMemoryUsage -= oldItem.Size
		delete(s.Data, key)
	}
}
