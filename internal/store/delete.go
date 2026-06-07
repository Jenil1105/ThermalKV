package store

import "fmt"

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
