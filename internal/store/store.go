package store

type Store struct {
	Data map[string]string
}

func NewStore() *Store {
	return &Store{
		Data: make(map[string]string),
	}
}
