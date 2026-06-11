package store

import "time"

func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock() // <-----
	item, exists := s.Data[key]

	if !exists {

		meltItem, exists := s.Thermal.LoadFromCool(key)

		if exists {

			if s.isExpired(meltItem) {

				s.Thermal.DeleteIndex(key)

				s.Mutex.RUnlock() // ----->
				return "", false
			}

			s.Mutex.RUnlock() // ----->

			s.Mutex.Lock() // </-----

			s.putItem(key, meltItem)

			if s.CurrentMemoryUsage > s.MaxHotMemory {
				go s.RunEmergencyCooling()
			}

			s.Thermal.DeleteIndex(key)

			s.Mutex.Unlock() // -----/>

			return meltItem.Value, true
		}

		s.Mutex.RUnlock() // ----->
		return "", false
	}

	// If item has expired, remove it and return as missing.
	if s.isExpired(item) {
		s.Mutex.RUnlock() // ----->

		s.Mutex.Lock() // </-----
		// Double-check under write lock before deleting to avoid races.
		item2, exists2 := s.Data[key]
		if exists2 {
			if s.isExpired(item2) {
				s.removeItem(key)
			}
		}
		s.Mutex.Unlock() // -----/>

		return "", false
	}

	value := item.Value
	s.Mutex.RUnlock()

	s.Mutex.Lock() // </-----
	item.LastAccessUnix = time.Now().Unix()
	s.Data[key] = item
	s.Mutex.Unlock() // -----/>

	return value, true
}
