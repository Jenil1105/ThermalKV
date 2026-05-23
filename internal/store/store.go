package store

import "time"

type Item struct {
	Value  string
	Expiry time.Time
}

type Store struct {
	Data map[string]Item
}

func NewStore() *Store {
	return &Store{
		Data: make(map[string]Item),
	}
}

func (s *Store) Set(key string, value string) {
	s.Data[key] = Item{
		Value: value,
	}
}

func (s *Store) Get(key string) (string, bool) {
	item, exists := s.Data[key]

	if !exists {
		return "", false
	}

	if !item.Expiry.IsZero() && time.Now().After(item.Expiry) {
		delete(s.Data, key)
		return "", false
	}

	return item.Value, true
}

func (s *Store) Delete(key string) {
	delete(s.Data, key)
}

func (s *Store) SetTTL(key string, seconds int) {
	item, exists := s.Data[key]

	if !exists {
		return
	}

	item.Expiry = time.Now().Add(time.Duration(seconds) * time.Second)

	s.Data[key] = item
}
