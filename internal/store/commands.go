package store

import (
	"fmt"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence"
	"time"
)

// Set
func (s *Store) Set(key string, value string) {

	s.Mutex.Lock()

	s.Data[key] = model.Item{
		Value: value,
	}

	s.WriteCount++
	s.Mutex.Unlock()
	s.WAL.Write("SET", key, value)

	if s.WriteCount >= 5 {
		snapshot := s.ExportData()
		persistence.SaveSnapshot(snapshot)
		persistence.ClearWAL()
		s.WriteCount = 0
		fmt.Println("Snapshot saved...")
	}

}

// Get
func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock()
	item, exists := s.Data[key]

	if !exists {
		s.Mutex.RUnlock()
		return "", false
	}

	// If item has expired, remove it and return as missing.
	if item.Expiry != 0 && time.Now().After(time.Unix(item.Expiry, 0)) {
		s.Mutex.RUnlock()
		s.Mutex.Lock()
		// Double-check under write lock before deleting to avoid races.
		item2, exists2 := s.Data[key]
		if exists2 {
			if item2.Expiry != 0 && time.Now().After(time.Unix(item2.Expiry, 0)) {
				delete(s.Data, key)
			}
		}
		s.Mutex.Unlock()
		return "", false
	}

	value := item.Value
	s.Mutex.RUnlock()
	return value, true
}

// Delete
func (s *Store) Delete(key string) {

	s.Mutex.Lock()

	delete(s.Data, key)
	s.WriteCount++
	s.Mutex.Unlock()
	s.WAL.Write("DEL", key)

	if s.WriteCount >= 5 {
		snapshot := s.ExportData()
		persistence.SaveSnapshot(snapshot)
		persistence.ClearWAL()
		s.WriteCount = 0
		fmt.Println("Snapshot saved...")
	}
}

// Count the total number of keys
func (s *Store) Count() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return len(s.Data)
}

// If key exists or not
func (s *Store) Exists(key string) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, exists := s.Data[key]

	return exists
}

// Get all key names
func (s *Store) Keys() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	keys := make([]string, 0, len(s.Data))
	for key := range s.Data {
		keys = append(keys, key)
	}
	return keys
}

// Get Heap size
func (s *Store) HeapSize() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return len(s.ExpiryHeap)
}
