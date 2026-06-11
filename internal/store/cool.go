package store

import "fmt"

func (s *Store) CoolKey(key string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item, exists := s.Data[key]

	if !exists {
		return fmt.Errorf("key not found")
	}
	err := s.Thermal.MoveToCool(key, item)

	if err != nil {
		return err
	}

	s.removeItem(key)

	return nil
}
