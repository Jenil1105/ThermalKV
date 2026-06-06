package store

import (
	"fmt"
	"thermalkv/internal/model"
	"time"
)

func (s *Store) StartCoolingWorker() {
	go func() {
		for {
			time.Sleep(500 * time.Second)

			var keysToCool []string

			s.Mutex.RLock()

			for key, item := range s.Data {
				if s.CoolingScore(item) > s.CoolingThreshold {
					keysToCool = append(keysToCool, key)
				}
			}

			s.Mutex.RUnlock()

			for _, key := range keysToCool {
				err := s.CoolKey(key)
				s.TotalCoolings++
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

	s.Mutex.Lock()

	if s.CoolingInProgress {
		s.Mutex.Unlock()
		return
	}

	s.CoolingInProgress = true
	s.Mutex.Unlock()

	defer func() {
		s.Mutex.Lock()
		s.CoolingInProgress = false
		s.Mutex.Unlock()
	}()

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
		s.TotalCoolings++

		if err != nil {
			return
		}

		fmt.Println("Emerfency cooled:", maxKey, "Score:", maxScore)

	}

}
