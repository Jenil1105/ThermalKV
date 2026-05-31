package store

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence"
	"thermalkv/internal/ttl"
	"time"
)

type Store struct {
	Data       map[string]model.Item
	Mutex      sync.RWMutex
	WriteCount int
	ExpiryHeap ttl.MinHeap
	WAL        *persistence.WAL
}

// NewStore
func NewStore(wal *persistence.WAL) *Store {
	h := ttl.MinHeap{}
	heap.Init(&h)
	return &Store{
		Data:       make(map[string]model.Item),
		ExpiryHeap: h,
		WAL:        wal,
	}
}

// Set
func (s *Store) Set(key string, value string) {

	s.Mutex.Lock()

	s.Data[key] = model.Item{
		Value: value,
	}

	s.WriteCount++
	s.Mutex.Unlock()
	s.WAL.Write("SET", key, value)

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
	s.WriteCount++
	s.Mutex.Unlock()
	s.WAL.Write("DEL", key)

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
	s.Data[key] = item

	heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
		Key:    key,
		Expiry: expiry,
	})

	s.WAL.Write("EXPIRE", key, strconv.FormatInt(expiry, 10))
}

// StartCleaner starts a background goroutine that periodically checks for expired keys and removes them from the store.
func (s *Store) StartCleaner() {
	go func() {
		for {
			s.Mutex.Lock()

			if len(s.ExpiryHeap) == 0 {
				s.Mutex.Unlock()
				time.Sleep(1 * time.Second)
				continue
			}

			top := s.ExpiryHeap[0]
			now := time.Now().Unix()

			if top.Expiry > now {
				sleepDuration := time.Duration(top.Expiry-now) * time.Second
				s.Mutex.Unlock()
				time.Sleep(sleepDuration)
				continue
			}

			expired := heap.Pop(&s.ExpiryHeap).(ttl.ExpiryItem)

			item, exists := s.Data[expired.Key]

			if exists {
				if item.Expiry == expired.Expiry {
					delete(s.Data, expired.Key)
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

		case "DEL":
			delete(s.Data, key)

		case "GET":
			// no recovery action needed

		case "EXPIRE":

			if len(parts) < 3 {
				continue
			}

			expiry, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				continue
			}
			item, exists := s.Data[key]
			if !exists {
				continue
			}
			s.Data[key] = item
			item.Expiry = expiry

			heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
				Key:    key,
				Expiry: expiry,
			})
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

		if item.Expiry != 0 && time.Now().Unix() >= item.Expiry {
			continue
		}

		s.Data[key] = model.Item{
			Value:  item.Value,
			Expiry: item.Expiry,
		}

		if item.Expiry != 0 {
			heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
				Key:    key,
				Expiry: item.Expiry,
			})
		}
	}
}

func (s *Store) Count() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return len(s.Data)
}

func (s *Store) Exists(key string) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, exists := s.Data[key]

	return exists
}

func (s *Store) Keys() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	keys := make([]string, 0, len(s.Data))
	for key := range s.Data {
		keys = append(keys, key)
	}
	return keys
}

func (s *Store) HeapSize() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return len(s.ExpiryHeap)
}

func (s *Store) GetWriteCount() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.WriteCount

}
