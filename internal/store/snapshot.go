package store

import (
	"thermalkv/internal/model"
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

		item := model.Item{
			Value:          item.Value,
			Expiry:         item.Expiry,
			Size:           int64(len(item.Value)),
			LastAccessUnix: time.Now().Unix(),
		}
		s.putItem(key, item)

		if item.Expiry != 0 {
			s.scheduleExpiry(key, item.Expiry)
		}
	}
}
