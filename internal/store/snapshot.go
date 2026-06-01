package store

import (
	"container/heap"
	"thermalkv/internal/model"
	"thermalkv/internal/ttl"
	"time"
)

func (s *Store) ExportData() map[string]model.SnapshotItem {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	snapshot := make(map[string]model.SnapshotItem)

	for key, item := range s.Data {
		snapshot[key] = model.SnapshotItem{
			Value:  item.Value,
			Expiry: item.Expiry,
		}
	}
	return snapshot
}

func (s *Store) ImportData(snapshot map[string]model.SnapshotItem) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for key, item := range snapshot {

		if item.Expiry != 0 && time.Now().Unix() >= item.Expiry {
			continue
		}

		s.Data[key] = model.Item{
			Value:  item.Value,
			Expiry: item.Expiry,
		}

		if item.Expiry != 0 {
			heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
				Key:    key,
				Expiry: item.Expiry,
			})
		}
	}
}
