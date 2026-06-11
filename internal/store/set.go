package store

import (
	"fmt"
	"thermalkv/internal/model"
	"time"
)

func (s *Store) Set(key string, value string) {

	size := int64(len(value))
	s.Mutex.Lock()

	item := model.Item{
		Value:          value,
		LastAccessUnix: time.Now().Unix(),
		Size:           size,
	}

	s.putItem(key, item)

	if s.CurrentMemoryUsage > s.MaxHotMemory {
		go s.RunEmergencyCooling()
	}

	s.Mutex.Unlock()
	err := s.WAL.Write("SET", key, value)

	if err != nil {
		fmt.Println("WAL write failed: ", err)
		return
	}

}
