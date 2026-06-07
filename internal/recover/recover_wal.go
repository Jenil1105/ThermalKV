package recover

import (
	"container/heap"
	"strconv"
	"strings"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/store"
	"thermalkv/internal/ttl"
	"time"
)

// Recover replays the given logs to restore the store's state after a restart.
func StoreLogs(s *store.Store, logs []string) {
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
			size := int64(len(value))
			s.Data[key] = model.Item{
				Value:          value,
				LastAccessUnix: time.Now().Unix(),
				Size:           size,
			}
			s.CurrentMemoryUsage += size

		case "DEL":

			if s.Exists(key) {
				s.CurrentMemoryUsage -= s.Data[key].Size
				delete(s.Data, key)
			}

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
			item.Expiry = expiry
			s.Data[key] = item

			heap.Push(&s.ExpiryHeap, ttl.ExpiryItem{
				Key:    key,
				Expiry: expiry,
			})
		}
	}
}

func RecoverWAL(s *store.Store, dir string) {

	logs := walpkg.LoadLogs(dir)
	StoreLogs(s, logs)

}
