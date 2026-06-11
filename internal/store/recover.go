package store

import (
	"thermalkv/internal/model"
	"time"
)

func (s *Store) RecoverSet(key, value string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item := model.Item{
		Value:          value,
		LastAccessUnix: time.Now().Unix(),
		Size:           int64(len(value)),
	}
	s.putItem(key, item)

}

func (s *Store) RecoverDelete(key string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if _, exists := s.Data[key]; exists {
		s.removeItem(key)
	}
}

func (s *Store) RecoverExpire(key string, expiry int64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	isExpirySet := s.setItemExpiry(key, expiry)
	if !isExpirySet {
		return
	}
}
