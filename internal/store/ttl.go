package store

import (
	"container/heap"
	"strconv"
	"thermalkv/internal/ttl"
	"time"
)

// Set TTL
func (s *Store) SetTTL(key string, seconds int) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item, exists := s.Data[key]

	if !exists {
		return
	}
	expiry := time.Now().Unix() + int64(seconds)
	item.Expiry = expiry
	s.Data[key] = item

	heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
		Key:    key,
		Expiry: expiry,
	})

	s.WAL.Write("EXPIRE", key, strconv.FormatInt(expiry, 10))
}
