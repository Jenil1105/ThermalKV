package store

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
