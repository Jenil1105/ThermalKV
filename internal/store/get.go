package store

import "time"

func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock() // <-----
	item, exists := s.Data[key]

	if !exists {

		item, exists = s.Thermal.LoadFromCool(key)

		if exists {

			if item.Expiry != 0 &&
				time.Now().Unix() > item.Expiry {

				s.Thermal.DeleteIndex(key)

				s.Mutex.RUnlock() // ----->
				return "", false
			}

			s.Mutex.RUnlock() // ----->

			s.Mutex.Lock() // </-----

			s.Data[key] = item
			s.CurrentMemoryUsage += item.Size
			if s.CurrentMemoryUsage > s.MaxHotMemory {
				go s.RunEmergencyCooling()
			}

			s.Thermal.DeleteIndex(key)

			s.Mutex.Unlock() // -----/>

			return item.Value, true
		}

		s.Mutex.RUnlock() // ----->
		return "", false
	}

	// If item has expired, remove it and return as missing.
	if item.Expiry != 0 && time.Now().Unix() > item.Expiry {
		s.Mutex.RUnlock() // ----->

		s.Mutex.Lock() // </-----
		// Double-check under write lock before deleting to avoid races.
		item2, exists2 := s.Data[key]
		if exists2 {
			if item2.Expiry != 0 && time.Now().Unix() > item2.Expiry {
				s.CurrentMemoryUsage -= item2.Size
				delete(s.Data, key)
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
