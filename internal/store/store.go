package store

import (
	"container/heap"
	"sync"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence"
	"thermalkv/internal/ttl"
	"time"
)

type Store struct {
	Data       map[string]model.Item
	Mutex      sync.RWMutex
	WriteCount int
	ExpiryHeap ttl.MinHeap
	WAL        *persistence.WAL
}

// NewStore
func NewStore(wal *persistence.WAL) *Store {
	h := ttl.MinHeap{}
	heap.Init(&h)
	return &Store{
		Data:       make(map[string]model.Item),
		ExpiryHeap: h,
		WAL:        wal,
	}
}

// StartCleaner starts a background goroutine that periodically checks for expired keys and removes them from the store.
func (s *Store) StartCleaner() {
	go func() {
		for {
			s.Mutex.Lock()

			if len(s.ExpiryHeap) == 0 {
				s.Mutex.Unlock()
				time.Sleep(1 * time.Second)
				continue
			}

			top := s.ExpiryHeap[0]
			now := time.Now().Unix()

			if top.Expiry > now {
				sleepDuration := time.Duration(top.Expiry-now) * time.Second
				s.Mutex.Unlock()
				time.Sleep(sleepDuration)
				continue
			}

			expired := heap.Pop(&s.ExpiryHeap).(ttl.ExpiryItem)

			item, exists := s.Data[expired.Key]

			if exists {
				if item.Expiry == expired.Expiry {
					delete(s.Data, expired.Key)
				}
			}
			s.Mutex.Unlock()
		}
	}()
}

func (s *Store) GetWriteCount() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.WriteCount

}
