package store

import (
	"container/heap"
	"strconv"
	"strings"
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

// Recover replays the given logs to restore the store's state after a restart.
func (s *Store) Recover(logs []string) {
	for _, log := range logs {
		parts := strings.Fields(log)

		if len(parts) < 2 {
			continue
		}

		operation := parts[0]
		key := parts[1]

		switch operation {

		case "SET":
			if len(parts) < 3 {
				continue
			}

			value := parts[2]
			s.Data[key] = model.Item{Value: value}

		case "DEL":
			delete(s.Data, key)

		case "GET":
			// no recovery action needed

		case "EXPIRE":

			if len(parts) < 3 {
				continue
			}

			expiry, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				continue
			}
			item, exists := s.Data[key]
			if !exists {
				continue
			}
			s.Data[key] = item
			item.Expiry = expiry

			heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
				Key:    key,
				Expiry: expiry,
			})
		}
	}
}

func (s *Store) GetWriteCount() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.WriteCount

}
