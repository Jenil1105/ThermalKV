package store

import (
	"fmt"
	"thermalkv/internal/persistence"
	"time"
)

func (s *Store) StartSnapshotLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			snapshot := s.ExportData()
			err := persistence.SaveSnapshot(snapshot)

			if err != nil {
				fmt.Println("Snapshout error: ", err)
				continue
			}

			fmt.Println("Snapshot saved")
		}
	}()
}
