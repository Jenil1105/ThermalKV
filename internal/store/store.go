package store

type Store struct {
	Data map[string]string
}

func NewStore() *Store {
	return &Store{
		Data: make(map[string]string),
	}
}

func (s *Store) Set(key string, value string) {
	s.Data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	value, exists := s.Data[key]
	return value, exists
}
