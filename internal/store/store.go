package store

import (
	"strings"
	"sync"
	"thermalkv/internal/persistence"
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
	persistence.WriteLog("SET", key, value)
}

func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock()
	item, exists := s.Data[key]

	if !exists {
		s.Mutex.RUnlock()
		return "", false
	}

	// If item has expired, remove it and return as missing.
	if !item.Expiry.IsZero() && time.Now().After(item.Expiry) {
		s.Mutex.RUnlock()
		s.Mutex.Lock()
		// Double-check under write lock before deleting to avoid races.
		item2, exists2 := s.Data[key]
		if exists2 {
			if !item2.Expiry.IsZero() && time.Now().After(item2.Expiry) {
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

func (s *Store) Recover(logs []string) {
	for _, log := range logs {
		parts := strings.Split(log, " ")

		if len(parts) < 3 {
			continue
		}

		operation := parts[0]
		key := parts[1]
		value := parts[2]

		if operation == "SET" {
			s.Data[key] = Item{Value: value}
		}
	}
}
