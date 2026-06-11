package store

import (
	"container/heap"
	"thermalkv/internal/model"
	"thermalkv/internal/ttl"
	"time"
)

func (s *Store) isExpired(item model.Item) bool {
	return item.Expiry != 0 && time.Now().Unix() >= item.Expiry
}

func (s *Store) scheduleExpiry(key string, expiry int64) {
	heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
		Key:    key,
		Expiry: expiry,
	})
}

func (s *Store) setItemExpiry(key string, expiry int64) bool {
	item, exists := s.Data[key]
	if !exists {
		return false
	}

	item.Expiry = expiry
	s.Data[key] = item
	s.scheduleExpiry(key, expiry)
	return true
}
