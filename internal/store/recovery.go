package store

import (
	"container/heap"
	"strconv"
	"strings"
	"thermalkv/internal/model"
	"thermalkv/internal/ttl"
)

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
