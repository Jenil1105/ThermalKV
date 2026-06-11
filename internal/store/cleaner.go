package store

import (
	"container/heap"
	"thermalkv/internal/ttl"
	"time"
)

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
					s.removeItem(expired.Key)
				}
			}
			s.Mutex.Unlock()
		}
	}()
}
