package store

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence"
	"time"
)

type Store struct {
	Data       map[string]model.Item
	Mutex      sync.RWMutex
	WriteCount int
}

// NewStore
func NewStore() *Store {
	return &Store{
		Data: make(map[string]model.Item),
	}
}

// Set
func (s *Store) Set(key string, value string) {

	s.Mutex.Lock()

	s.Data[key] = model.Item{
		Value: value,
	}
	persistence.WriteLog("SET", key, value)

	s.Mutex.Unlock()

	s.WriteCount++

	if s.WriteCount >= 5 {
		snapshot := s.ExportData()
		persistence.SaveSnapshot(snapshot)
		persistence.ClearWAL()
		s.WriteCount = 0
		fmt.Println("Snapshot saved...")
	}

}

// Get
func (s *Store) Get(key string) (string, bool) {

	s.Mutex.RLock()
	item, exists := s.Data[key]

	if !exists {
		s.Mutex.RUnlock()
		return "", false
	}

	// If item has expired, remove it and return as missing.
	if item.Expiry != 0 && time.Now().After(time.Unix(item.Expiry, 0)) {
		s.Mutex.RUnlock()
		s.Mutex.Lock()
		// Double-check under write lock before deleting to avoid races.
		item2, exists2 := s.Data[key]
		if exists2 {
			if item2.Expiry != 0 && time.Now().After(time.Unix(item2.Expiry, 0)) {
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

// Delete
func (s *Store) Delete(key string) {

	s.Mutex.Lock()

	delete(s.Data, key)
	persistence.WriteLog("DEL", key)

	s.Mutex.Unlock()

	s.WriteCount++
	if s.WriteCount >= 5 {
		snapshot := s.ExportData()
		persistence.SaveSnapshot(snapshot)
		persistence.ClearWAL()
		s.WriteCount = 0
		fmt.Println("Snapshot saved...")
	}
}

func (s *Store) SetTTL(key string, seconds int) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item, exists := s.Data[key]

	if !exists {
		return
	}
	expiry := time.Now().Unix() + int64(seconds)
	item.Expiry = expiry
	persistence.WriteLog("EXPIRE", key, strconv.FormatInt(expiry, 10))
	s.Data[key] = item
}

// StartCleaner starts a background goroutine that periodically checks for expired keys and removes them from the store.
func (s *Store) StartCleaner() {
	go func() {
		for {
			time.Sleep(1 * time.Second)

			s.Mutex.Lock()

			for key, item := range s.Data {
				if item.Expiry != 0 && time.Now().After(time.Unix(item.Expiry, 0)) {
					delete(s.Data, key)
				}
			}

			s.Mutex.Unlock()
		}
	}()
}

// Recover replays the given logs to restore the store's state after a restart.
func (s *Store) Recover(logs []string) {
	for _, log := range logs {
		parts := strings.Fields(log)

		if len(parts) < 2 {
			continue
		}

		operation := parts[0]
		key := parts[1]

		switch operation {

		case "SET":
			if len(parts) < 3 {
				continue
			}

			value := parts[2]
			s.Data[key] = model.Item{Value: value}

		case "DELETE":
			delete(s.Data, key)

		case "GET":
			// no recovery action needed

		case "EXPIRE":
			expiry, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				continue
			}
			item, exists := s.Data[key]
			if exists {
				item.Expiry = expiry
				s.Data[key] = item
			} else {
				return
			}
		}
	}
}

func (s *Store) ExportData() map[string]model.SnapshotItem {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	snapshot := make(map[string]model.SnapshotItem)

	for key, item := range s.Data {
		snapshot[key] = model.SnapshotItem{
			Value:  item.Value,
			Expiry: item.Expiry,
		}
	}
	return snapshot
}

func (s *Store) ImportData(snapshot map[string]model.SnapshotItem) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for key, item := range snapshot {
		s.Data[key] = model.Item{
			Value:  item.Value,
			Expiry: item.Expiry,
		}
	}
}
