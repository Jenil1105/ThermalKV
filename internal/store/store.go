package store

import (
	"container/heap"
	"fmt"
	"sync"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence"
	"thermalkv/internal/thermal"
	"thermalkv/internal/ttl"
	"time"
)

type Store struct {
	Data               map[string]model.Item
	CurrentMemoryUsage int64
	MaxHotMemory       int64
	CoolingThreshold   int64
	Mutex              sync.RWMutex
	ExpiryHeap         ttl.MinHeap
	WAL                *persistence.WAL
	Thermal            *thermal.Manager
}

// NewStore
func NewStore(wal *persistence.WAL, manager *thermal.Manager) *Store {
	h := ttl.MinHeap{}
	heap.Init(&h)
	return &Store{
		Data:             make(map[string]model.Item),
		MaxHotMemory:     100,
		CoolingThreshold: 100000,
		ExpiryHeap:       h,
		WAL:              wal,
		Thermal:          manager,
	}
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

func (s *Store) CoolKey(key string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	item, exists := s.Data[key]
	s.CurrentMemoryUsage -= item.Size

	if !exists {
		return fmt.Errorf("key not found")
	}
	err := s.Thermal.MoveToCool(key, item)

	if err != nil {
		return err
	}

	delete(s.Data, key)

	return nil
}

func (s *Store) StartCoolingWorker() {
	go func() {
		for {
			time.Sleep(500 * time.Second)

			now := time.Now().Unix()

			var keysToCool []string

			s.Mutex.RLock()

			for key, item := range s.Data {
				if now-item.LastAccessUnix > 60 {
					keysToCool = append(keysToCool, key)
				}
			}

			s.Mutex.RUnlock()

			for _, key := range keysToCool {
				err := s.CoolKey(key)

				if err == nil {
					fmt.Println("Auto cooled:", key)
				}
			}

		}
	}()
}

func (s *Store) CoolingScore(
	item model.Item,
) int64 {

	idleTime := time.Now().Unix() - item.LastAccessUnix

	return item.Size * idleTime
}

func (s *Store) RunEmergencyCooling() {

	for {

		if s.CurrentMemoryUsage <= s.MaxHotMemory {
			return
		}

		var maxKey string
		var maxScore int64

		s.Mutex.RLock()

		for key, item := range s.Data {
			score := s.CoolingScore(item)

			if score > maxScore {
				maxScore = score
				maxKey = key
			}
		}

		s.Mutex.RUnlock()

		if maxKey == "" {
			return
		}

		err := s.CoolKey(maxKey)

		if err != nil {
			return
		}

		fmt.Println("Emerfency cooled:", maxKey, "Score:", maxScore)

	}

}
