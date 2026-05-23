package store

import (
	"sync"
	"time"
)

type Item struct {
	Value  string
	Expiry time.Time
}

type Store struct {
	Data  map[string]Item
	Mutex sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		Data: make(map[string]Item),
	}
}

func (s *Store) Set(key string, value string) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Data[key] = Item{
		Value: value,
	}
}

func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	item, exists := s.Data[key]

	if !exists {
		return "", false
	}

	return item.Value, true
}

func (s *Store) Delete(key string) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	delete(s.Data, key)
}

func (s *Store) SetTTL(key string, seconds int) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item, exists := s.Data[key]

	if !exists {
		return
	}

	item.Expiry = time.Now().Add(time.Duration(seconds) * time.Second)

	s.Data[key] = item
}

func (s *Store) StartCleaner() {
	go func() {
		for {
			time.Sleep(1 * time.Second)

			s.Mutex.Lock()

			for key, item := range s.Data {
				if !item.Expiry.IsZero() && time.Now().After(item.Expiry) {
					println("we need to delete key ", key)
					delete(s.Data, key)
				}
			}

			s.Mutex.Unlock()
		}
	}()
}
