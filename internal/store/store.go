package store

import (
	"container/heap"
	"fmt"
	"os"
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
	CoolingInProgress  bool
	TotalPromotions    int64
	TotalCompactions   int64
	TotalCoolings      int64
}

// NewStore
func NewStore(wal *persistence.WAL, manager *thermal.Manager) *Store {
	h := ttl.MinHeap{}
	heap.Init(&h)
	return &Store{
		Data:              make(map[string]model.Item),
		MaxHotMemory:      100,
		CoolingThreshold:  100000,
		ExpiryHeap:        h,
		WAL:               wal,
		Thermal:           manager,
		CoolingInProgress: false,
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

func (s *Store) GetInfo() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	info := []string{
		"===== ThermalKV Info =====",
		fmt.Sprintf(
			"HOT Keys           : %d",
			len(s.Data),
		),
		fmt.Sprintf(
			"COOL Keys          : %d",
			len(s.Thermal.ColdIndex),
		),
		fmt.Sprintf(
			"HOT Memory Usage   : %d bytes",
			s.CurrentMemoryUsage,
		),
		fmt.Sprintf(
			"Max HOT Memory     : %d bytes",
			s.MaxHotMemory,
		),
		fmt.Sprintf(
			"Cooling Threshold  : %d",
			s.CoolingThreshold,
		),
		fmt.Sprintf(
			"Cold File Size     : %d",
			GetColdFileSize(),
		),
	}

	return info

}

func GetColdFileSize() int64 {
	fileInfo, err := os.Stat("data/cold.dat")

	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
