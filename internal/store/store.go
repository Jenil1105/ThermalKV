package store

import (
	"container/heap"
	"fmt"
	"sync"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence"
	"thermalkv/internal/thermal"
	"thermalkv/internal/ttl"
	"time"
)

type Store struct {
	Data       map[string]model.Item
	Mutex      sync.RWMutex
	WriteCount int
	ExpiryHeap ttl.MinHeap
	WAL        *persistence.WAL
	Thermal    *thermal.Manager
}

// NewStore
func NewStore(wal *persistence.WAL, manager *thermal.Manager) *Store {
	h := ttl.MinHeap{}
	heap.Init(&h)
	return &Store{
		Data:       make(map[string]model.Item),
		ExpiryHeap: h,
		WAL:        wal,
		Thermal:    manager,
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

func (s *Store) CoolKey(key string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item, exists := s.Data[key]

	if !exists {
		return fmt.Errorf("key not found")
	}
	err := s.Thermal.MoveToCool(key, item)

	if err != nil {
		return err
	}

	delete(s.Data, key)

	return nil
}
