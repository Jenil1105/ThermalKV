package store

import (
	"fmt"
	"thermalkv/internal/model"
	"time"
)

// Set
func (s *Store) Set(key string, value string) {

	size := int64(len(value))
	s.Mutex.Lock()

	oldItem, exists := s.Data[key]

	if exists {
		s.CurrentMemoryUsage -= oldItem.Size
	}

	s.Data[key] = model.Item{
		Value:          value,
		LastAccessUnix: time.Now().Unix(),
		Size:           size,
	}
	s.CurrentMemoryUsage += size

	if s.CurrentMemoryUsage > s.MaxHotMemory {
		go s.RunEmergencyCooling()
	}

	s.Mutex.Unlock()
	err := s.WAL.Write("SET", key, value)

	if err != nil {
		fmt.Println("WAL write failed: ", err)
		return
	}

}

// Get
func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock()
	item, exists := s.Data[key]

	if !exists {

		item, exists = s.Thermal.LoadFromCool(key)

		if exists {

			if item.Expiry != 0 &&
				time.Now().Unix() > item.Expiry {

				delete(
					s.Thermal.ColdIndex,
					key,
				)

				s.Mutex.RUnlock()
				return "", false
			}

			s.Mutex.RUnlock()

			s.Mutex.Lock()

			s.Data[key] = item
			s.CurrentMemoryUsage += item.Size
			if s.CurrentMemoryUsage > s.MaxHotMemory {
				go s.RunEmergencyCooling()
			}

			delete(
				s.Thermal.ColdIndex,
				key,
			)

			s.Mutex.Unlock()

			return item.Value, true
		}

		s.Mutex.RUnlock()
		return "", false
	}

	item.LastAccessUnix = time.Now().Unix()
	s.Data[key] = item

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

	if s.Exists(key) {
		s.Mutex.Lock()
		delete(s.Data, key)
		s.Mutex.Unlock()
	} else {
		_, exists := s.Thermal.ColdIndex[key]

		if exists {
			err := s.Thermal.AppendDelete(key)
			if err != nil {
				return
			}
			delete(s.Thermal.ColdIndex, key)
			return
		}
	}

	err := s.WAL.Write("DEL", key)
	if err != nil {
		fmt.Println("WAL write failed: ", err)
		return
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
func (s *Store) HotMemory() int64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.CurrentMemoryUsage
}
